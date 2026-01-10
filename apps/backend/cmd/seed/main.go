package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/projectrole"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	skipSeed := flag.Bool("skip-seed", false, "Skip seeding data, only reset schema and apply migrations")
	flag.Parse()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			logrus.Fatalf("Failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	// Hard reset: drop and recreate the public schema
	rawDB, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("failed opening raw db connection: %v", err)
	}
	if _, err := rawDB.ExecContext(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public; DROP SCHEMA IF EXISTS atlas_schema_revisions CASCADE;`); err != nil {
		logrus.Fatalf("failed resetting schema: %v", err)
	}
	if err := rawDB.Close(); err != nil {
		logrus.Fatalf("failed closing raw db: %v", err)
	}
	logrus.Info("Database schema reset (dropped and recreated).")

	logrus.Info("Applying database migrations...")
	cmd := exec.Command("pnpm", "run", "db:migrate:apply")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("failed running database migrations: %v", err)
	}
	logrus.Info("Database migrations applied.")

	if *skipSeed {
		logrus.Info("Skipping data seeding as requested.")
		return
	}

	// Default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("failed to hash password: %v", err)
	}

	// Seed maps
	personIDs := make(map[string]uuid.UUID)
	orgNodeIDs := make(map[string]uuid.UUID)
	productIDs := make(map[string]uuid.UUID) // key: DOI

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
		SetIsSysAdmin(true).
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		logrus.Fatalf("failed creating user for %s: %v", testUserName, err)
	}
	logrus.Infof("Created user for person %s", testUserName)

	es := eventstore.NewEntStore(client)

	// --- Seed Roles / Scopes / Memberships for org tree ---

	orgRepo := organisationrepo.NewEntRepo(client)
	personRepo := personrepo.NewEntRepo(client)
	orgSvc := organisation.NewService(orgRepo, personRepo)

	orgRoot, err := orgSvc.CreateRoot(ctx, "Nederland", nil, nil, nil)
	if err != nil {
		logrus.Fatalf("create root org node: %v", err)
	}
	orgNodeIDs["Nederland"] = orgRoot.ID

	roleRepo := projectrole.NewRepository(client)

	if _, err := roleRepo.Create(ctx, "contributor", "Contributor", orgRoot.ID); err != nil {
		logrus.Fatalf("create project role contributor: %v", err)
	}
	if _, err := roleRepo.Create(ctx, "admin", "Project Lead", orgRoot.ID); err != nil {
		logrus.Fatalf("create project role admin: %v", err)
	}

	projects := []seedProject{
		{
			Title:        "Quantum-Resistant Cryptography Benchmarking",
			Description:  "Evaluating performance and security of post-quantum algorithms across diverse architectures.",
			Organisation: "Cybersecurity Lab – Utrecht University",
			People:       []string{"Dr. Elaine Carter", "Tomas Ternovski", "Prof. Jin-Ho Park"},
			Start:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			End:          time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC),
			Products: []seedProduct{
				{Type: entities.Software, Language: "en", Name: "PQCryptoBench", DOI: "10.1234/pqcb.2024.001"},
				{Type: entities.Dataset, Language: "en", Name: "Post-Quantum Benchmark Dataset", DOI: "10.1234/pqcb.2024.002"},
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
				{Type: entities.Dataset, Language: "en", Name: "Methane Emission Field Measurements", DOI: "10.1234/mmc.2024.001"},
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
				{Type: entities.Software, Language: "en", Name: "AIBench", DOI: "10.1234/alam.2024.001"},
				{Type: entities.Dataset, Language: "en", Name: "AIBench Dataset", DOI: "10.1234/alam.2024.002"},
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
				{Type: entities.Software, Language: "en", Name: "WaveSoft", DOI: "10.1234/wbhp.2024.001"},
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
				{Type: entities.Software, Language: "en", Name: "Marine Drone Swarms", DOI: "10.1234/mdsm.2024.001"},
			},
		},
	}

	// Ensure each project includes the test user
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

	// Helpers
	mustPersonID := func(name string) uuid.UUID {
		id, ok := personIDs[name]
		if !ok {
			logrus.Fatalf("no person ID found for %q", name)
		}
		return id
	}
	mustOrgNodeID := func(name string) uuid.UUID {
		id, ok := orgNodeIDs[name]
		if !ok {
			logrus.Fatalf("no org node ID found for %q", name)
		}
		return id
	}
	mustProductID := func(doi string) uuid.UUID {
		id, ok := productIDs[doi]
		if !ok {
			logrus.Fatalf("no product ID found for DOI %q", doi)
		}
		return id
	}

	universities, err := getOrCreateChild(ctx, client, orgSvc, orgRoot.ID, "Universities")
	if err != nil {
		logrus.Fatalf("create/get Universities node: %v", err)
	}

	uuLeafID, err := createPath(ctx, orgSvc, universities.ID,
		"Utrecht University",
		"Faculty of Science",
		"Department of Information and Computing Sciences",
		"Cybersecurity Group",
		"Post-Quantum Cryptography Lab",
	)
	if err != nil {
		logrus.Fatalf("seed UU subtree: %v", err)
	}
	orgNodeIDs["Cybersecurity Lab – Utrecht University"] = uuLeafID

	tudLeafID, err := createPath(ctx, orgSvc, universities.ID,
		"TU Delft",
		"Faculty of Electrical Engineering, Mathematics and Computer Science",
		"Department of Intelligent Systems",
		"Distributed Graphics Group",
		"Wave Rendering Lab",
	)
	if err != nil {
		logrus.Fatalf("seed TU Delft subtree: %v", err)
	}
	orgNodeIDs["Distributed Graphics Lab – TU Delft"] = tudLeafID

	researchInstitutes, err := getOrCreateChild(ctx, client, orgSvc, orgRoot.ID, "Research Institutes")
	if err != nil {
		logrus.Fatalf("create/get Institutes node: %v", err)
	}

	medaiLeafID, err := createPath(ctx, orgSvc, researchInstitutes.ID,
		"MedTech & AI",
		"MedAI Institute Rotterdam",
		"Clinical Decision Systems",
		"Adaptive Diagnostics Unit",
	)
	if err != nil {
		logrus.Fatalf("seed MedAI subtree: %v", err)
	}
	orgNodeIDs["MedAI Institute Rotterdam"] = medaiLeafID

	agroLeafID, err := createPath(ctx, orgSvc, orgRoot.ID,
		"Applied Research",
		"Agri & Food",
		"AgroTech Consortium",
		"AgroTech Research Group",
		"Greenhouse Emissions Program",
	)
	if err != nil {
		logrus.Fatalf("seed AgroTech subtree: %v", err)
	}
	orgNodeIDs["AgroTech Research Group"] = agroLeafID

	oceanLeafID, err := createPath(ctx, orgSvc, researchInstitutes.ID,
		"Ocean & Robotics",
		"Ocean Robotics Centre",
		"Leiden Site",
		"Autonomous Swarms Division",
	)
	if err != nil {
		logrus.Fatalf("seed Ocean subtree: %v", err)
	}
	orgNodeIDs["Ocean Robotics Centre Leiden"] = oceanLeafID

	for _, sp := range projects {
		// People
		var authorIDs []uuid.UUID
		for _, name := range sp.People {
			if _, exists := personIDs[name]; exists {
				authorIDs = append(authorIDs, personIDs[name])
				continue
			}

			userID := uuid.New()
			email := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, ".", ""), " ", ".")) + "@example.com"

			p, err := client.Person.
				Create().
				SetName(name).
				SetUserID(userID).
				SetAvatarURL(avatarUrl).
				SetDescription(defaultBio).
				SetEmail(email).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating person %q: %v", name, err)
			}

			personIDs[name] = p.ID
			authorIDs = append(authorIDs, p.ID)
			logrus.Infof("Created person %s (%s)", name, p.ID)

			_, err = client.User.
				Create().
				SetID(userID).
				SetPersonID(p.ID).
				SetPassword(string(hashedPassword)).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating user for person %q: %v", name, err)
			}
			logrus.Infof("Created user for person %s", name)
		}

		// Products (create once; reuse IDs later in event stream)
		for _, prod := range sp.Products {
			if _, exists := productIDs[prod.DOI]; exists {
				continue
			}

			row, err := client.Product.
				Create().
				SetName(prod.Name).
				SetType(int(prod.Type)).
				SetLanguage(prod.Language).
				SetDoi(prod.DOI).
				AddAuthorIDs(authorIDs...).
				Save(ctx)
			if err != nil {
				logrus.Fatalf("failed creating product %q: %v", prod.Name, err)
			}

			productIDs[prod.DOI] = row.ID
			logrus.Infof("Created product %s (%s)", prod.Name, row.ID)
		}
	}

	// --- Seed Roles / Scopes / Memberships for org tree ---

	// Helper to create roles for an org
	createRolesForOrg := func(orgID uuid.UUID) (adminRoleID, researcherRoleID, studentsRoleID uuid.UUID) {
		orgEnt, err := client.OrganisationNode.Get(ctx, orgID)
		if err != nil {
			logrus.Fatalf("failed getting org %s for roles: %v", orgID, err)
		}

		// Admin
		adminRole, err := client.OrganisationRole.Create().
			SetKey("admin").
			SetDisplayName("Administrator").
			SetOrganisation(orgEnt).
			SetPermissions([]string{
				"manage_members",
				"manage_project_roles",
				"manage_organisation_roles",
				"manage_custom_fields",
				"manage_details",
			}).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create admin role for %s: %v", orgID, err)
		}

		// Create legacy Scope for Admin
		adminScope, err := client.RoleScope.Create().
			SetRole(adminRole).
			SetRootNode(orgEnt).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create admin scope for %s: %v", orgID, err)
		}

		// Researcher
		researcherRole, err := client.OrganisationRole.Create().
			SetKey("researcher").
			SetDisplayName("Researcher").
			SetOrganisation(orgEnt).
			SetPermissions([]string{}). // Basic access
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create researcher role for %s: %v", orgID, err)
		}

		researcherScope, err := client.RoleScope.Create().
			SetRole(researcherRole).
			SetRootNode(orgEnt).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create researcher scope: %v", err)
		}

		// Students
		studentsRole, err := client.OrganisationRole.Create().
			SetKey("students").
			SetDisplayName("Student").
			SetOrganisation(orgEnt).
			SetPermissions([]string{}).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create students role for %s: %v", orgID, err)
		}

		studentsScope, err := client.RoleScope.Create().
			SetRole(studentsRole).
			SetRootNode(orgEnt).
			Save(ctx)
		if err != nil {
			logrus.Fatalf("create students scope: %v", err)
		}

		return adminScope.ID, researcherScope.ID, studentsScope.ID
	}

	// Create roles for Root Org
	rootAdminScopeID, rootResearcherScopeID, rootStudentsScopeID := createRolesForOrg(orgRoot.ID)

	// For demo: create specific roles for sub-orgs if needed, or re-use root roles?
	// The previous seed had "Students" scoped to a subtree. Now roles are strictly per-org.
	// If "Students" was scoped to "Cybersecurity Lab", we must create a role there.

	// studentsRootID := orgRoot.ID (removed unused)

	studentsScopeID := rootStudentsScopeID

	if id, ok := orgNodeIDs["Cybersecurity Lab – Utrecht University"]; ok {
		// studentsRootID = id // removed unused
		// Create roles for this sub-org
		_, _, subStudentsScopeID := createRolesForOrg(id)
		studentsScopeID = subStudentsScopeID
	}

	// Memberships: example assignments
	_, err = client.Membership.Create().SetPersonID(mustPersonID(testUserName)).SetRoleScopeID(rootAdminScopeID).Save(ctx)
	if err != nil {
		logrus.Fatalf("create admin membership: %v", err)
	}

	if _, ok := personIDs["Dr. Elaine Carter"]; ok {
		_, err = client.Membership.Create().SetPersonID(mustPersonID("Dr. Elaine Carter")).SetRoleScopeID(rootResearcherScopeID).Save(ctx)
		if err != nil {
			logrus.Fatalf("create researcher membership: %v", err)
		}
	}
	if _, ok := personIDs["Tomas Ternovski"]; ok {
		// Use the students scope we determined earlier (Root or Subtree)
		_, err = client.Membership.Create().SetPersonID(mustPersonID("Tomas Ternovski")).SetRoleScopeID(studentsScopeID).Save(ctx)
		if err != nil {
			logrus.Fatalf("create students membership: %v", err)
		}
	}

	logrus.Info("Seeding projects (event stream)...")

	for _, sp := range projects {
		projectID := uuid.New()

		startEvent := &events.ProjectStarted{
			Base: events.Base{
				ID:        uuid.New(),
				ProjectID: projectID,
				At:        time.Now().UTC(),
				Status:    "approved",
			},
			Title:           sp.Title,
			Description:     sp.Description,
			StartDate:       sp.Start,
			EndDate:         sp.End,
			OwningOrgNodeID: mustOrgNodeID(sp.Organisation),
		}

		if err := es.Append(ctx, projectID, 0, startEvent); err != nil {
			logrus.Fatalf("append ProjectStarted for %s: %v", sp.Title, err)
		}

		version := 1

		// fetch project role IDs
		contributorRole, err := roleRepo.GetByKeyAndOrg(ctx, "contributor", orgRoot.ID)
		if err != nil {
			logrus.Fatalf("fetch contributor role: %v", err)
		}
		leadRole, err := roleRepo.GetByKeyAndOrg(ctx, "admin", orgRoot.ID)
		if err != nil {
			logrus.Fatalf("fetch admin role: %v", err)
		}

		contributorRoleID := contributorRole.ID
		leadRoleID := leadRole.ID

		for _, name := range sp.People {
			personID := mustPersonID(name)

			roleID := contributorRoleID
			if name == testUserName {
				roleID = leadRoleID // make test user the project lead/admin
			}

			pevt := &events.ProjectRoleAssigned{
				Base: events.Base{
					ID:        uuid.New(),
					ProjectID: projectID,
					At:        time.Now().UTC(),
					Status:    "approved",
				},
				PersonID:      personID,
				ProjectRoleID: roleID,
			}

			if err := es.Append(ctx, projectID, version, pevt); err != nil {
				logrus.Fatalf("append ProjectRoleAssigned for %s (%s): %v", name, sp.Title, err)
			}
			version++
		}

		for _, prod := range sp.Products {
			productID := mustProductID(prod.DOI)

			pevt := &events.ProductAdded{
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
		}

		logrus.Infof("Seeded project: %s (%s)", sp.Title, projectID.String())
	}

	logrus.Info("Seeding done.")

	// Notifications
	logrus.Info("Seeding notifications...")
	notificationRecipients := []string{
		"Dr. Elaine Carter",
		"Sarah Vos",
		"Dr. Mariam Bensaïd",
		"Niels van Bruggen",
		"Dr. Yara Mendes",
		"Emilio Vargas",
		testUserName,
	}

	for _, name := range notificationRecipients {
		personID, ok := personIDs[name]
		if !ok {
			continue
		}

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
			SetMessage("Your project has been started.").
			SetRead(true).
			SetSentAt(time.Now().Add(-48 * time.Hour)).
			Save(ctx)
		if err != nil {
			logrus.Errorf("failed to create notification for %s: %v", name, err)
		}
	}
	logrus.Info("Notifications seeded.")
}

// createPath creates a path of OrganisationNodes under the given root.
func createPath(ctx context.Context, orgSvc organisation.Service, rootID uuid.UUID, names ...string) (uuid.UUID, error) {
	parentID := rootID
	for _, name := range names {
		n, err := orgSvc.CreateChild(ctx, parentID, name, nil, nil, nil)
		if err != nil {
			return uuid.Nil, err
		}
		parentID = n.ID
	}
	return parentID, nil
}

// getOrCreateChild retrieves a child OrganisationNode by name under the given parent,
func getOrCreateChild(ctx context.Context, cli *ent.Client, orgSvc organisation.Service, parentID uuid.UUID, name string) (*entities.OrganisationNode, error) {
	row, err := cli.OrganisationNode.
		Query().
		Where(
			organisationnode.NameEQ(name),
			organisationnode.ParentIDEQ(parentID),
		).
		Only(ctx)
	if err == nil {
		return transform.ToEntityPtr[entities.OrganisationNode](row), nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}
	return orgSvc.CreateChild(ctx, parentID, name, nil, nil, nil)
}
