//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Load environment variables from .env file (backend or root)
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env") // Also try root .env

	ctx := context.Background()

	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("ent/migrate/migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating atlas migration directory")
	}

	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                          // provide migration directory
		schema.WithMigrationMode(schema.ModeInspect), // inspect existing database state
		schema.WithDialect(dialect.Postgres),         // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}

	if len(os.Args) < 2 {
		log.Fatal().Msg("migration name is required. Use: 'pnpm run db:migrate:diff <name>'")
	}

	migrationName := os.Args[1]

	// Build the database URL from environment variables
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Generate migrations using Atlas support for PostgreSQL
	err = migrate.NamedDiff(ctx, dbURL, migrationName, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed generating migration file")
	}

	log.Info().Msgf("Migration '%s' generated successfully in ent/migrate/migrations/", migrationName)
}
