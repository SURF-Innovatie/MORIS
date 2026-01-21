package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	email := flag.String("email", "", "Email of the admin user")
	name := flag.String("name", "", "Name of the admin user")
	password := flag.String("password", "", "Password for the admin user")
	flag.Parse()

	if *email == "" || *name == "" || *password == "" {
		log.Fatal().Msg("email, name, and password are required")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Global.DBHost, env.Global.DBPort, env.Global.DBUser, env.Global.DBPassword, env.Global.DBName)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to postgres")
	}
	defer client.Close()

	ctx := context.Background()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to hash password")
	}

	// Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start transaction")
	}

	// Create Person
	userAccountID := uuid.New()
	person, err := tx.Person.
		Create().
		SetName(*name).
		SetUserID(userAccountID).
		SetEmail(*email).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		log.Fatal().Err(err).Msg("failed creating person")
	}

	// Create User
	_, err = tx.User.
		Create().
		SetID(userAccountID).
		SetPersonID(person.ID).
		SetIsSysAdmin(true).
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		log.Fatal().Err(err).Msg("failed creating user")
	}

	if err := tx.Commit(); err != nil {
		log.Fatal().Err(err).Msg("failed to commit transaction")
	}

	log.Info().Msgf("Successfully created sys_admin user: %s (%s)", *name, *email)
}
