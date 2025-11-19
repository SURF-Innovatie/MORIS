package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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

	ctx := context.Background()
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logrus.Fatalf("failed running Ent database migrations: %v", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword123"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("failed to hash password: %v", err)
	}

	// Create test user
	testUser, err := client.User.
		Create().
		SetName("Test User").
		SetEmail("test@example.com").
		SetPassword(string(hashedPassword)).
		Save(ctx)

	if err != nil {
		logrus.Fatalf("failed creating test user: %v", err)
	}

	logrus.Infof("Successfully created test user with ID: %d", testUser.ID)
	logrus.Infof("Email: %s", testUser.Email)
	logrus.Infof("Name: %s", testUser.Name)
	logrus.Info("Password: testpassword123")
}
