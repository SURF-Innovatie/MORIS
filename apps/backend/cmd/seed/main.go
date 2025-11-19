package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type seedProject struct {
	Title        string
	Description  string
	Organisation string
	People       []string
	Start        time.Time
	End          time.Time
}

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
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			logrus.Fatalf("Failed to close client")
		}
	}(client)

	ctx := context.Background()

	// drop database and run migrations
	if err := client.Schema.Create(
		ctx,
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		logrus.Fatalf("failed running Ent database migrations: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword123"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("failed to hash password: %v", err)
	}

	client.User.Delete().ExecX(ctx)
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

	es := eventstore.NewEntStore(client)

	projects := []seedProject{
		{
			Title:        "Quantum-Resistant Cryptography Benchmarking",
			Description:  "Evaluating performance and security of post-quantum algorithms across diverse architectures.",
			Organisation: "Cybersecurity Lab – Utrecht University",
			People:       []string{"Dr. Elaine Carter", "Tomas Ternovski", "Prof. Jin-Ho Park"},
			Start:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:        "Microbial Methane Capture for Sustainable Farms",
			Description:  "Engineering microbial systems that reduce methane emission in agricultural environments.",
			Organisation: "AgroTech Research Group",
			People:       []string{"Sarah Vos", "Dr. Pieter de Louw", "Emilio Vargas"},
			Start:        time.Date(2024, 3, 12, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:        "Adaptive Learning Algorithms for Medical Diagnostics",
			Description:  "Developing adaptive neural decision systems for clinical diagnostics.",
			Organisation: "MedAI Institute Rotterdam",
			People:       []string{"Dr. Mariam Bensaïd", "Konrad Schulz", "Olivia Becker"},
			Start:        time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:        "Wave-Based Holographic Rendering on Edge Devices",
			Description:  "Investigating real-time holographic rendering techniques for small form-factor devices.",
			Organisation: "Distributed Graphics Lab – TU Delft",
			People:       []string{"Niels van Bruggen", "Prof. Hiro Tanaka", "Emily Rhodes"},
			Start:        time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:        "Marine Drone Swarms for Microplastic Detection",
			Description:  "Deploying autonomous micro-drones to map microplastic concentration gradients.",
			Organisation: "Ocean Robotics Centre Leiden",
			People:       []string{"Dr. Yara Mendes", "Stef Kranenburg", "Akira Watanabe"},
			Start:        time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 7, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	personIDs := make(map[string]uuid.UUID)

	for _, sp := range projects {
		for _, name := range sp.People {
			if _, exists := personIDs[name]; exists {
				continue
			}

			row, err := client.Person.
				Create().
				SetName(name).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating person %q: %v", name, err)
			}

			personIDs[name] = row.ID
			logrus.Infof("Created person %s (%s)", name, row.ID)
		}
	}

	logrus.Info("Seeding projects...")

	for _, sp := range projects {
		projectID := uuid.New()

		startEvent := events.ProjectStarted{
			Base: events.Base{
				ProjectID: projectID,
				At:        time.Now().UTC(),
			},
			Title:       sp.Title,
			Description: sp.Description,
			StartDate:   sp.Start,
			EndDate:     sp.End,
			Organisation: entities.Organisation{
				Id:   uuid.New(),
				Name: sp.Organisation,
			},
		}

		if err := es.Append(ctx, projectID, 0, startEvent); err != nil {
			logrus.Fatalf("append ProjectStarted for %s: %v", sp.Title, err)
		}

		version := 1
		for _, name := range sp.People {
			personID, ok := personIDs[name]
			if !ok {
				logrus.Fatalf("no person ID found for %q", name)
			}

			pevt := events.PersonAdded{
				Base: events.Base{
					ProjectID: projectID,
					At:        time.Now().UTC(),
				},
				PersonId: personID,
			}

			if err := es.Append(ctx, projectID, version, pevt); err != nil {
				logrus.Fatalf("append PersonAdded for %s (%s): %v", name, sp.Title, err)
			}
			version++
		}

		logrus.Infof("Seeded project: %s (%s)", sp.Title, projectID.String())
	}

	logrus.Info("Seeding done.")
}
