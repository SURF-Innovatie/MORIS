package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("Error loading .env file: %v", err)
	}

	email := flag.String("email", "", "Email of the admin user")
	name := flag.String("name", "", "Name of the admin user")
	password := flag.String("password", "", "Password for the admin user")
	flag.Parse()

	if *email == "" || *name == "" || *password == "" {
		logrus.Fatal("email, name, and password are required")
	}

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

	ctx := context.Background()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("failed to hash password: %v", err)
	}

	// Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		logrus.Fatalf("failed to start transaction: %v", err)
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
		logrus.Fatalf("failed creating person: %v", err)
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
		logrus.Fatalf("failed creating user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		logrus.Fatalf("failed to commit transaction: %v", err)
	}

	logrus.Infof("Successfully created sys_admin user: %s (%s)", *name, *email)
}
