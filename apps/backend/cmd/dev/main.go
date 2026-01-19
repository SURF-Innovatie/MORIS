package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	_ "github.com/SURF-Innovatie/MORIS/api/swag-docs"
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/api"
	"github.com/SURF-Innovatie/MORIS/internal/di"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
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
	// Run Atlas migrations
	log.Info().Msg("Running database migrations...")

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
		log.Fatal().Err(err).Msg("failed applying migrations")
	}
	log.Info().Msg("Database migrations applied successfully")

	// Validate event registrations
	if err := events.ValidateRegistrations(); err != nil {
		log.Fatal().Err(err).Msg("event registration invalid")
	}

	injector := do.New(di.Package)
	defer injector.Shutdown() //nolint:errcheck

	// Get ent client to defer close
	client := do.MustInvoke[*ent.Client](injector)
	defer client.Close()

	r := api.SetupRouter(injector)

	log.Info().Msgf("Go Backend Server starting on http://localhost:%s", env.Global.Port)
	log.Fatal().Err(http.ListenAndServe(":"+env.Global.Port, r)).Msg("Server failed")
}
