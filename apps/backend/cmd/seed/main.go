package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/migrate"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type seedProject struct {
	Title        string
	Description  string
	Organisation string
	People       []string
	Products     []seedProduct
	Start        time.Time
	End          time.Time
}

type seedProduct struct {
	Type        entities.ProductType
	Language    string
	Name        string
	DOI         string
	AuthorNames []string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warnf("Error loading .env file: %v", err)
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
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			logrus.Fatalf("Failed to close client")
		}
	}(client)

	ctx := context.Background()

	// Hard reset: drop and recreate the public schema
	rawDB, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("failed opening raw db connection: %v", err)
	}
	if _, err := rawDB.ExecContext(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`); err != nil {
		logrus.Fatalf("failed resetting schema: %v", err)
	}
	if err := rawDB.Close(); err != nil {
		logrus.Fatalf("failed closing raw db: %v", err)
	}
	logrus.Info("Database schema reset (dropped and recreated).")

	// drop database and run migrations
	if err := client.Schema.Create(
		ctx,
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logrus.Fatalf("failed running Ent database migrations: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("failed to hash password: %v", err)
	}

	personIDs := make(map[string]uuid.UUID)
	organisationIDs := make(map[string]uuid.UUID)

	const testUserName = "Test User"
	const testUserEmail = "test.user@example.com"
	const avatarUrl = "https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp&f=y"
	const defaultBio = "This is a test user account created during seeding."
	testUserAccountID := uuid.New()

	testPerson, err := client.Person.
		Create().
		SetName(testUserName).
		SetUserID(testUserAccountID).
		SetEmail(testUserEmail).
		SetAvatarURL(avatarUrl).
		SetDescription(defaultBio).
		Save(ctx)
	if err != nil {
		logrus.Fatalf("failed creating %s person: %v", testUserName, err)
	}
	testPersonID := testPerson.ID
	personIDs[testUserName] = testPersonID
	logrus.Infof("Created person %s (%s)", testUserName, testPersonID)

	_, err = client.User.
		Create().
		SetID(testUserAccountID).
		SetPersonID(testPersonID).
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		logrus.Fatalf("failed creating user for %s: %v", testUserName, err)
	}
	logrus.Infof("Created user for person %s", testUserName)

	es := eventstore.NewEntStore(client)

	projects := []seedProject{
		{
			Title:        "Quantum-Resistant Cryptography Benchmarking",
			Description:  "Evaluating performance and security of post-quantum algorithms across diverse architectures.",
			Organisation: "Cybersecurity Lab – Utrecht University",
			People:       []string{"Dr. Elaine Carter", "Tomas Ternovski", "Prof. Jin-Ho Park"},
			Start:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{
					Type:     entities.Software,
					Language: "en",
					Name:     "PQCryptoBench",
					DOI:      "10.1234/pqcb.2024.001",
				},
				{
					Type:     entities.Dataset,
					Language: "en",
					Name:     "Post-Quantum Benchmark Dataset",
					DOI:      "10.1234/pqcb.2024.002",
				},
			},
		},
		{
			Title:        "Microbial Methane Capture for Sustainable Farms",
			Description:  "Engineering microbial systems that reduce methane emission in agricultural environments.",
			Organisation: "AgroTech Research Group",
			People:       []string{"Emilio Vargas", "Sarah Vos", "Dr. Pieter de Louw"},
			Start:        time.Date(2024, 3, 12, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{
					Type:     entities.Dataset,
					Language: "en",
					Name:     "Methane Emission Field Measurements",
					DOI:      "10.1234/mmc.2024.001",
				},
			},
		},
		{
			Title:        "Adaptive Learning Algorithms for Medical Diagnostics",
			Description:  "Developing adaptive neural decision systems for clinical diagnostics.",
			Organisation: "MedAI Institute Rotterdam",
			People:       []string{"Dr. Mariam Bensaïd", "Konrad Schulz", "Olivia Becker"},
			Start:        time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{
					Type:     entities.Software,
					Language: "en",
					Name:     "AIBench",
					DOI:      "10.1234/alam.2024.001",
				},
				{
					Type:     entities.Dataset,
					Language: "en",
					Name:     "AIBench Dataset",
					DOI:      "10.1234/alam.2024.002",
				},
			},
		},
		{
			Title:        "Wave-Based Holographic Rendering on Edge Devices",
			Description:  "Investigating real-time holographic rendering techniques for small form-factor devices.",
			Organisation: "Distributed Graphics Lab – TU Delft",
			People:       []string{"Niels van Bruggen", "Prof. Hiro Tanaka", "Emily Rhodes"},
			Start:        time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{
					Type:     entities.Software,
					Language: "en",
					Name:     "WaveSoft",
					DOI:      "10.1234/wbhp.2024.001",
				},
			},
		},
		{
			Title:        "Marine Drone Swarms for Microplastic Detection",
			Description:  "Deploying autonomous micro-drones to map microplastic concentration gradients.",
			Organisation: "Ocean Robotics Centre Leiden",
			People:       []string{"Dr. Yara Mendes", "Stef Kranenburg", "Akira Watanabe"},
			Start:        time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 7, 12, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{
					Type:     entities.Software,
					Language: "en",
					Name:     "Marine Drone Swarms",
					DOI:      "10.1234/mdsm.2024.001",
				},
			},
		},
	}

	for i := range projects {
		hasTestUser := false
		for _, person := range projects[i].People {
			if person == testUserName {
				hasTestUser = true
				break
			}
		}
		if !hasTestUser {
			projects[i].People = append(projects[i].People, testUserName)
		}
	}

	for _, sp := range projects {
		var authorIds []uuid.UUID

		for _, name := range sp.People {
			if _, exists := personIDs[name]; exists {
				authorIds = append(authorIds, personIDs[name])
				continue
			}

			userID := uuid.New()

			row, err := client.Person.
				Create().
				SetName(name).
				SetUserID(userID).
				SetAvatarURL(avatarUrl).
				SetDescription(defaultBio).
				SetEmail(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, ".", ""), " ", ".")) + "@example.com").
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating person %q: %v", name, err)
			}
			authorIds = append(authorIds, row.ID)

			personIDs[name] = row.ID
			logrus.Infof("Created person %s (%s)", name, row.ID)

			// We make each person a user with a default password
			_, err = client.User.
				Create().
				SetID(userID).
				SetPersonID(row.ID).
				SetPassword(string(hashedPassword)).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating user for person %q: %v", name, err)
			}
			logrus.Infof("Created user for person %s", name)
		}

		for _, prod := range sp.Products {
			row, err := client.Product.
				Create().
				SetName(prod.Name).
				SetType(int(prod.Type)).
				SetLanguage(prod.Language).
				SetDoi(prod.DOI).
				AddAuthorIDs(authorIds...).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating product %q: %v", prod.Name, err)
			}

			logrus.Infof("Created product %s (%s)", prod.Name, row.ID)
		}

		org := sp.Organisation

		if _, exists := organisationIDs[org]; exists {
			continue
		}

		row, err := client.Organisation.
			Create().
			SetName(org).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("failed creating organisation %q: %v", org, err)
		}

		organisationIDs[org] = row.ID
		logrus.Infof("Created organisation %s (%s)", org, row.ID)
	}

	logrus.Info("Seeding projects...")

	for _, sp := range projects {
		projectID := uuid.New()

		startEvent := events.ProjectStarted{
			Base: events.Base{
				ID:        uuid.New(),
				ProjectID: projectID,
				At:        time.Now().UTC(),
				Status:    "approved",
			},
			ProjectAdmin:   testPersonID,
			Title:          sp.Title,
			Description:    sp.Description,
			StartDate:      sp.Start,
			EndDate:        sp.End,
			OrganisationID: organisationIDs[sp.Organisation],
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
					Status:    "approved",
				},
				PersonId: personID,
			}

			if err := es.Append(ctx, projectID, version, pevt); err != nil {
				logrus.Fatalf("append PersonAdded for %s (%s): %v", name, sp.Title, err)
			}
			version++
		}

		for _, prod := range sp.Products {
			productID := uuid.New()

			row, err := client.Product.
				Create().
				SetID(productID).
				SetName(prod.Name).
				SetType(int(prod.Type)).
				SetLanguage(prod.Language).
				SetDoi(prod.DOI).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating product %q: %v", prod.Name, err)
			}
			logrus.Infof("Created product %s (%s)", row.Name, row.ID)

			pevt := events.ProductAdded{
				Base: events.Base{
					ProjectID: projectID,
					At:        time.Now().UTC(),
					Status:    "approved",
				},
				ProductID: productID,
			}

			if err := es.Append(ctx, projectID, version, pevt); err != nil {
				logrus.Fatalf("append ProductAdded for %s (%s): %v", prod.Name, sp.Title, err)
			}
			version++

			logrus.Infof("Seeded product %s (%s) for project %s", prod.Name, productID, sp.Title)
		}
		logrus.Infof("Seeded project: %s (%s)", sp.Title, projectID.String())
	}

	logrus.Info("Seeding done.")

	// Add some sample notifications
	logrus.Info("Seeding notifications...")
	notificationRecipients := []string{"Dr. Elaine Carter", "Sarah Vos", "Dr. Mariam Bensaïd", "Niels van Bruggen", "Dr. Yara Mendes", "Emilio Vargas", testUserName}
	for _, name := range notificationRecipients {
		personID, ok := personIDs[name]
		if !ok {
			continue
		}
		// Find user ID for this person
		u, err := client.User.Query().Where(entuser.PersonIDEQ(personID)).Only(ctx)
		if err != nil {
			logrus.Errorf("failed to find user for person %s: %v", name, err)
			continue
		}

		_, err = client.Notification.Create().
			SetUser(u).
			SetMessage("Welcome to MORIS! This is a sample notification.").
			SetRead(false).
			SetSentAt(time.Now().Add(-24 * time.Hour)).
			Save(ctx)
		if err != nil {
			logrus.Errorf("failed to create notification for %s: %v", name, err)
		}

		_, err = client.Notification.Create().
			SetUser(u).
			SetMessage("Your project 'Quantum-Resistant Cryptography Benchmarking' has been started.").
			SetRead(true).
			SetSentAt(time.Now().Add(-48 * time.Hour)).
			Save(ctx)
		if err != nil {
			logrus.Errorf("failed to create notification for %s: %v", name, err)
		}
	}
	logrus.Info("Notifications seeded.")
}
