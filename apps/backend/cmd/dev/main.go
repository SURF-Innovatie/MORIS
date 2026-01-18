package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	exzenodo "github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	surfconextapp "github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/app/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	customfieldrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/customfield"
	errorlogrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/error_log"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/SURF-Innovatie/MORIS/api/swag-docs"
	"github.com/SURF-Innovatie/MORIS/ent"
	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
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
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events/hydrator"
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
	zenodohandler "github.com/SURF-Innovatie/MORIS/internal/handler/zenodo"
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

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	raidsink "github.com/SURF-Innovatie/MORIS/internal/adapter/sinks/raid"
	csvsource "github.com/SURF-Innovatie/MORIS/internal/adapter/sources/csv"
	adapterhandler "github.com/SURF-Innovatie/MORIS/internal/handler/adapter"

	eventpolicysvc "github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	eventpolicyhandler "github.com/SURF-Innovatie/MORIS/internal/handler/eventpolicy"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
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

	// Auto-migration is disabled in favor of versioned migrations (Atlas)
	// if err := client.Schema.Create(
	// 	context.Background(),
	// Run Atlas migrations
	logrus.Info("Running database migrations...")

	// Construct URL safely to handle special characters in password
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(env.Global.DBUser, env.Global.DBPassword),
		Host:     fmt.Sprintf("%s:%s", env.Global.DBHost, env.Global.DBPort),
		Path:     env.Global.DBName,
		RawQuery: "sslmode=disable",
	}

	cmd := exec.Command("atlas", "migrate", "apply",
		"--url", u.String(),
		"--dir", "file://ent/migrate/migrations",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("failed applying migrations: %v", err)
	}
	logrus.Info("Database migrations applied successfully")

	// Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.Global.CacheHost, env.Global.CachePort),
		Password: env.Global.CachePassword,
		Username: env.Global.CacheUser,
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

	orcidOpts := env.ORCIDOptionsFromEnv()

	orcidCli := exorcid.NewClient(http.DefaultClient, orcidOpts)
	orcidSvc := orcid.NewService(userRepo, personRepo, orcidCli)

	surfOpts := env.SurfconextOptionsFromEnv()
	surfClient := exsurfconext.NewClient(http.DefaultClient, surfOpts)
	surfSvc := surfconextapp.NewService(surfClient, authSvc)
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
	customFieldRepo := customfieldrepo.NewRepository(client)
	customFieldSvc := customfield.NewService(customFieldRepo)

	organisationSvc := organisation.NewService(orgRepo, personRepo, rbacSvc)
	organisationHandler := organisationhandler.NewHandler(organisationSvc, rbacSvc, roleSvc, customFieldSvc)

	crossrefConfig := &excrossref.Config{
		BaseURL:   "https://api.crossref.org",
		UserAgent: "MORIS/1.0 (mailto:support@moris.org)",
		Mailto:    "support@moris.org",
	}

	crossrefClient := excrossref.NewClient(crossrefConfig)
	crossrefSvc := crossref.NewService(crossrefClient)
	crossrefHandler := crossrefhandler.NewHandler(crossrefSvc)

	notifRepo := notificationrepo.NewEntRepo(client)
	notifierSvc := notification.NewService(notifRepo)

	errorLogRepo := errorlogrepo.NewRepository(client)
	errorLogSvc := errorlog.NewService(errorLogRepo)

	eventSvc := event.NewService(esStore, client, notifierSvc)

	eventSvc.RegisterNotificationHandler(&event.ProjectEventNotificationHandler{Cli: client, ES: esStore})
	eventSvc.RegisterNotificationHandler(&event.ApprovalRequestNotificationHandler{Cli: client, ES: esStore, RBAC: rbacSvc})

	notificationHandler := notificationhandler.NewHandler(notifierSvc)

	// Create HTTP handler/controller
	authHandler := authhandler.NewHandler(userSvc, authSvc, orcidSvc, surfSvc)
	orcidHandler := orcidhandler.NewHandler(orcidSvc)

	zenodoOpts := env.ZenodoOptionsFromEnv()
	zenodoClient := exzenodo.NewClient(http.DefaultClient, zenodoOpts)
	zenodoSvc := zenodo.NewService(userRepo, zenodoClient)
	zenodoHandler := zenodohandler.NewHandler(zenodoSvc, curUser)

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

	projSvc := queries.NewService(esStore, ldr, repo, roleRepo, curUser, userSvc)
	projHandler := projecthandler.NewHandler(projSvc, customFieldSvc)

	evtHydrator := hydrator.New(repo, repo, repo, repo, userSvc)
	evtHandler := eventHandler.NewHandler(eventSvc, projSvc, userSvc, client, evtHydrator)

	statusUpdateHandler := &event.StatusUpdateNotificationHandler{
		Cli:            client,
		Hydrator:       evtHydrator,
		ProjectService: projSvc,
	}
	eventSvc.RegisterStatusChangeHandler(statusUpdateHandler.Handle)

	// Event Policies
	eventPolicyRepo := eventpolicyrepo.NewEntRepository(client)
	orgClosureProvider := eventpolicyrepo.NewOrgClosureAdapter(orgRepo)
	eventPolicySvc := eventpolicysvc.NewService(eventPolicyRepo, orgClosureProvider)
	eventPolicyHandler := eventpolicyhandler.NewHandler(eventPolicySvc)

	// Policy Evaluator Components
	recipientResolver := eventpolicyrepo.NewRecipientAdapter(client)
	notificationAdapter := eventpolicyrepo.NewNotificationAdapter(client)
	policyEvaluator := eventpolicy.NewEvaluator(eventPolicyRepo, orgClosureProvider, recipientResolver, notificationAdapter, evtHydrator)

	projCmdSvc := command.NewService(esStore, eventSvc, cacheSvc, refreshSvc, curUser, entProv, roleSvc, policyEvaluator, organisationSvc, rbacSvc)
	projCmdHandler := commandHandler.NewHandler(projCmdSvc)

	eventSvc.RegisterNotificationHandler(&event.EventPolicyHandler{PolicyRepo: eventPolicyRepo, Cli: client})
	eventSvc.RegisterNotificationHandler(&event.PolicyExecutionHandler{Evaluator: policyEvaluator, ProjectSvc: projSvc})

	userHandler := userhandler.NewHandler(userSvc, projSvc)

	// Adapters (Sources/Sinks)
	registry := adapter.NewRegistry()
	registry.RegisterSource(csvsource.NewCSVSource("/tmp/import.csv"))

	// RAiD configuration
	raidOpts := raid.DefaultOptions()
	// raidOpts.Username = env.Global.RAiDUsername // TODO: Add to env
	// raidOpts.Password = env.Global.RAiDPassword // TODO: Add to env
	raidClient := raid.NewClient(http.DefaultClient, raidOpts)

	registry.RegisterSink(raidsink.NewRAiDSink(raidClient))
	adapterHandler := adapterhandler.NewHandler(registry, projSvc)

	// Router
	r := chi.NewRouter()
	log := logrus.New()
	r.Use(logger.Logger("router", log))
	r.Use(middleware.Recoverer)

	r.Route(env.Global.APIBasePath, func(r chi.Router) {
		r.Use(authmiddleware.ErrorLoggingMiddleware(errorLogSvc))
		authhandler.MountRoutes(r, authSvc, authHandler)
		systemhandler.MountRoutes(r, systemHandler)
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware(authSvc))
			r.Route("/projects", func(r chi.Router) {
				projecthandler.MountProjectRoutes(r, projHandler)
				r.Get("/{id}/roles", projHandler.ListAvailableRoles)
				commandHandler.MountProjectCommandRouter(r, projCmdHandler)
				// Event policy routes for projects
				eventPolicyHandler.RegisterProjectRoutes(r)
			})
			organisationhandler.MountOrganisationRoutes(r, organisationHandler, rbacHandler)
			eventHandler.MountEventRoutes(r, evtHandler)
			personhandler.MountPersonRoutes(r, personHandler)
			orcidhandler.MountRoutes(r, orcidHandler)
			zenodohandler.MountRoutes(r, zenodoHandler)
			producthandler.MountProductRoutes(r, productHandler)
			notificationhandler.MountNotificationRoutes(r, notificationHandler)
			userhandler.MountUserRoutes(r, userHandler)
			crossrefhandler.MountCrossrefRoutes(r, crossrefHandler)
			adapterhandler.MountRoutes(r, adapterHandler)

			// Event Policies routes (standalone and org-scoped)
			eventPolicyHandler.RegisterRoutes(r)
			r.Route("/organisations/{id}/policies", func(r chi.Router) {
				r.Get("/", eventPolicyHandler.ListForOrgNode)
				r.Post("/", eventPolicyHandler.CreateForOrgNode)
			})
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
		httpSwagger.URL("swagger.json"), // the url pointing to API definition
	))

	logrus.Fatal(http.ListenAndServe(":"+env.Global.Port, r))
}
