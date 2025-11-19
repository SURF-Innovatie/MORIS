package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/handler/custom"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/project"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
	userSvc := user.NewService(client)
	authSvc := auth.NewService(client)
	
	// Set auth service for middleware
	auth.SetAuthService(authSvc)

	// Create HTTP handler/controller
	customHandler := custom.NewHandler(userSvc, authSvc)

	esStore := eventstore.NewEntStore(client)
	projSvc := project.NewService(esStore, client)
	projHandler := projecthandler.NewHandler(projSvc)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		custom.MountCustomHandlers(r, customHandler)
		projecthandler.MountProjectRoutes(r, projHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("$PORT must be set")
	}
	logrus.Infof("Go Backend Server starting on http://localhost:%s", port)
	logrus.Fatal(http.ListenAndServe(":"+port, r))
}
