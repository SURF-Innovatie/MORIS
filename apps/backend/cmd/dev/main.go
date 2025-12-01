package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	_ "github.com/SURF-Innovatie/MORIS/api/swag-docs"
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/handler/custom"
	notificationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/notification"
	organisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/organisation"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	producthandler "github.com/SURF-Innovatie/MORIS/internal/handler/product"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	userhandler "github.com/SURF-Innovatie/MORIS/internal/handler/user"
	notification "github.com/SURF-Innovatie/MORIS/internal/notification"
	"github.com/SURF-Innovatie/MORIS/internal/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/person"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/product"
	"github.com/SURF-Innovatie/MORIS/internal/project"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

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

	// Create services
	personSvc := person.NewService(client)
	userSvc := user.NewService(client, personSvc)
	authSvc := auth.NewService(client, userSvc)
	orcidSvc := orcid.NewService(client, userSvc)

	personHandler := personhandler.NewHandler(personSvc)
	productSvc := product.NewService(client)
	productHandler := producthandler.NewHandler(productSvc)

	organisationSvc := organisation.NewService(client)
	organisationHandler := organisationhandler.NewHandler(organisationSvc)

	notifierSvc := notification.NewService(client)

	// Set auth service for middleware
	auth.SetAuthService(authSvc)

	// Create HTTP handler/controller
	customHandler := custom.NewHandler(userSvc, authSvc, orcidSvc)

	userHandler := userhandler.NewHandler(userSvc)

	esStore := eventstore.NewEntStore(client)
	projSvc := project.NewService(esStore, client, notifierSvc)
	projHandler := projecthandler.NewHandler(projSvc)

	notificationHandler := notificationhandler.NewHandler(notifierSvc)

	// Router
	r := chi.NewRouter()
	log := logrus.New()
	r.Use(logger.Logger("router", log))
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		custom.MountCustomHandlers(r, customHandler)
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware)
			projecthandler.MountProjectRoutes(r, projHandler)
			projecthandler.MountEventRoutes(r, projHandler)
			personhandler.MountPersonRoutes(r, personHandler)
			producthandler.MountProductRoutes(r, productHandler)
			organisationhandler.MountOrganisationRoutes(r, organisationHandler)
			notificationhandler.MountNotificationRoutes(r, notificationHandler)
			userhandler.MountUserRoutes(r, userHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("$PORT must be set")
	}
	logrus.Infof("Go Backend Server starting on http://localhost:%s", port)
	// Serve the generated swagger JSON and assets and the Swagger UI at /swagger/
	r.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swag-docs/swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/swagger.json"), // the url pointing to API definition
	))

	logrus.Fatal(http.ListenAndServe(":"+port, r))
}
