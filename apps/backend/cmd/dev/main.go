package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/SURF-Innovatie/MORIS/api/swag-docs"
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	crossref2 "github.com/SURF-Innovatie/MORIS/external/crossref"
	"github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	appproduct "github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/cachewarmup"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	projectrolesvc "github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/env"
	"github.com/SURF-Innovatie/MORIS/internal/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	authhandler "github.com/SURF-Innovatie/MORIS/internal/handler/auth"
	crossrefhandler "github.com/SURF-Innovatie/MORIS/internal/handler/crossref"
	eventHandler "github.com/SURF-Innovatie/MORIS/internal/handler/event"
	authmiddleware "github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	notificationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/notification"
	orcidhandler "github.com/SURF-Innovatie/MORIS/internal/handler/orcid"
	organisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/organisation"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	producthandler "github.com/SURF-Innovatie/MORIS/internal/handler/product"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	commandHandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/command"
	systemhandler "github.com/SURF-Innovatie/MORIS/internal/handler/system"
	userhandler "github.com/SURF-Innovatie/MORIS/internal/handler/user"
	"github.com/SURF-Innovatie/MORIS/internal/infra/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/entclient"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	notificationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/notification"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	organisationrbacrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation_rbac"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	productrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/product"
	projectmembershiprepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_membership"
	projectquery "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project_query"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/projectrole"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
)

// @title MORIS
// @version 1.0
// @description MORIS
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" then a space and your JWT token.

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Global.DBHost, env.Global.DBPort, env.Global.DBUser, env.Global.DBPassword, env.Global.DBName)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logrus.Fatalf("failed running Ent database migrations: %v", err)
	}

	// Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.Global.CacheHost, env.Global.CachePort),
		Password: env.Global.CachePassword,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logrus.Warnf("failed to connect to redis/valkey: %v", err)
	} else {
		logrus.Infof("Connected to Redis at %s:%s", env.Global.CacheHost, env.Global.CachePort)
	}
	defer rdb.Close()

	if err := events.ValidateRegistrations(); err != nil {
		logrus.Fatalf("event registration invalid: %v", err)
	}

	personRepo := personrepo.NewEntRepo(client)

	userRepo := userrepo.NewEntRepo(client)
	membershipRepo := projectmembershiprepo.NewEntRepo(client)

	// Create services
	esStore := eventstore.NewEntStore(client)
	personSvc := personsvc.NewService(personRepo)
	userSvc := user.NewService(userRepo, personSvc, esStore, membershipRepo)
	authSvc := auth.NewJWTService(client, userSvc, personSvc, env.Global.JWTSecret)
	orcidSvc := orcid.NewService(client, userSvc, http.DefaultClient)
	curUser := auth.NewCurrentUserProvider(client)

	personHandler := personhandler.NewHandler(personSvc)
	productRepo := productrepo.NewEntRepo(client)
	productSvc := appproduct.NewService(productRepo)
	productHandler := producthandler.NewHandler(productSvc, curUser)

	orgRepo := organisationrepo.NewEntRepo(client)
	rbacRepo := organisationrbacrepo.NewEntRepo(client)

	rbacSvc := organisationrbac.NewService(rbacRepo)
	rbacHandler := organisationhandler.NewRBACHandler(rbacSvc)

	roleRepo := projectrole.NewRepository(client)
	roleSvc := projectrolesvc.NewService(roleRepo, orgRepo)
	customFieldSvc := customfield.NewService(client)

	organisationSvc := organisation.NewService(orgRepo, personRepo)
	organisationHandler := organisationhandler.NewHandler(organisationSvc, rbacSvc, roleSvc, customFieldSvc)

	crossrefConfig := &crossref2.Config{
		BaseURL:   "https://api.crossref.org",
		UserAgent: "MORIS/1.0 (mailto:support@moris.org)",
		Mailto:    "support@moris.org",
	}
	crossrefSvc := crossref2.NewService(crossrefConfig)
	crossrefHandler := crossrefhandler.NewHandler(crossrefSvc)

	notifRepo := notificationrepo.NewEntRepo(client)
	notifierSvc := notification.NewService(notifRepo)
	errorLogSvc := errorlog.NewService(client)

	eventSvc := event.NewService(esStore, client, notifierSvc)

	eventSvc.RegisterNotificationHandler(&event.ProjectEventNotificationHandler{Cli: client, ES: esStore})
	eventSvc.RegisterNotificationHandler(&event.ApprovalRequestNotificationHandler{Cli: client, ES: esStore, RBAC: rbacSvc})
	eventSvc.RegisterNotificationHandler(&event.StatusUpdateNotificationHandler{Cli: client})
	evtHandler := eventHandler.NewHandler(eventSvc, client)

	notificationHandler := notificationhandler.NewHandler(notifierSvc)


	// Create HTTP handler/controller
	authHandler := authhandler.NewHandler(userSvc, authSvc, orcidSvc)
	orcidHandler := orcidhandler.NewHandler(orcidSvc)
	systemHandler := systemhandler.NewHandler()

	cacheSvc := cache.NewRedisProjectCache(rdb, 24*time.Hour)
	refreshSvc := cache.NewEventstoreProjectCacheRefresher(esStore, cacheSvc)

	eventSvc.RegisterStatusChangeHandler(func(ctx context.Context, e events.Event) error {
		_, err := refreshSvc.Refresh(ctx, e.AggregateID())
		return err
	})

	repo := projectquery.NewEntRepo(client)
	ldr := load.New(esStore, cacheSvc)
	warmup := cachewarmup.NewService(repo, ldr, cacheSvc)
	entProv := entclient.New(client)


	projSvc := queries.NewService(esStore, ldr, repo, roleRepo, curUser)
	projHandler := projecthandler.NewHandler(projSvc, customFieldSvc)

	projCmdSvc := command.NewService(esStore, eventSvc, cacheSvc, refreshSvc, curUser, entProv)
	projCmdHandler := commandHandler.NewHandler(projCmdSvc)

	userHandler := userhandler.NewHandler(userSvc, projSvc)

	// Router
	r := chi.NewRouter()
	log := logrus.New()
	r.Use(logger.Logger("router", log))
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Use(authmiddleware.ErrorLoggingMiddleware(errorLogSvc))
		authhandler.MountRoutes(r, authSvc, authHandler)
		systemhandler.MountRoutes(r, systemHandler)
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware(authSvc))
			r.Route("/projects", func(r chi.Router) {
				projecthandler.MountProjectRoutes(r, projHandler)
				r.Get("/{id}/roles", projHandler.ListAvailableRoles)
				commandHandler.MountProjectCommandRouter(r, projCmdHandler)
			})
			organisationhandler.MountOrganisationRoutes(r, organisationHandler, rbacHandler)
			eventHandler.MountEventRoutes(r, evtHandler)
			personhandler.MountPersonRoutes(r, personHandler)
			orcidhandler.MountRoutes(r, orcidHandler)
			producthandler.MountProductRoutes(r, productHandler)
			notificationhandler.MountNotificationRoutes(r, notificationHandler)
			userhandler.MountUserRoutes(r, userHandler)
			crossrefhandler.MountCrossrefRoutes(r, crossrefHandler)
		})
	})


	// Warmup cache in background
	go func() {
		cached, err := warmup.WarmupProjects(context.Background())
		if err != nil {
			logrus.Errorf("Failed to warmup cache: %v", err)
		} else {
			logrus.Infof("Warmed up cache for %d projects", cached)
		}
	}()

	logrus.Infof("Go Backend Server starting on http://localhost:%s", env.Global.Port)
	// Serve the generated swagger JSON and assets and the Swagger UI at /swagger/
	r.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swag-docs/swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+env.Global.Port+"/swagger/swagger.json"), // the url pointing to API definition
	))

	logrus.Fatal(http.ListenAndServe(":"+env.Global.Port, r))
}
