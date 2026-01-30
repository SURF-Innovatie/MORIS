package role_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/role"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

func NewRoleClient(t *testing.T) *ent.Client {
	t.Helper()
	return enttest.Open(t, "sqlite3", "file:rbacrepo_test?mode=memory&cache=shared&_fk=1")
}

func seedOrgTree(t *testing.T, cli *ent.Client) (rootID, childID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	root, err := cli.OrganisationNode.Create().SetName("root").Save(ctx)
	if err != nil {
		t.Fatalf("create root: %v", err)
	}
	child, err := cli.OrganisationNode.Create().SetName("child").SetParentID(root.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create child: %v", err)
	}

	// closures (including self)
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(root.ID).SetDescendantID(root.ID).SetDepth(0).Save(ctx)
	if err != nil {
		t.Fatalf("closure root->root: %v", err)
	}
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(child.ID).SetDescendantID(child.ID).SetDepth(0).Save(ctx)
	if err != nil {
		t.Fatalf("closure child->child: %v", err)
	}
	_, err = cli.OrganisationNodeClosure.Create().SetAncestorID(root.ID).SetDescendantID(child.ID).SetDepth(1).Save(ctx)
	if err != nil {
		t.Fatalf("closure root->child: %v", err)
	}

	return root.ID, child.ID
}

func seedPersonUser(t *testing.T, cli *ent.Client, isSysAdmin bool) (personID, userID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	p, err := cli.Person.Create().
		SetName("Ada Lovelace").
		SetGivenName("Ada").
		SetFamilyName("Lovelace").
		SetEmail(uuid.NewString() + "@example.org").
		Save(ctx)
	if err != nil {
		t.Fatalf("create person: %v", err)
	}

	u, err := cli.User.Create().
		SetPersonID(p.ID).
		SetIsSysAdmin(isSysAdmin).
		Save(ctx)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return p.ID, u.ID
}

func TestListRoles_FiltersAndOrders(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)

	rootID, childID := seedOrgTree(t, cli)

	// Create two roles in root (unordered display names)
	_, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("b").
		SetDisplayName("Zeta").
		SetPermissions([]string{string(entities.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	_, err = cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("a").
		SetDisplayName("Alpha").
		SetPermissions([]string{string(entities.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	// Another role in child org
	_, err = cli.OrganisationRole.Create().
		SetOrganisationNodeID(childID).
		SetKey("c").
		SetDisplayName("Other").
		SetPermissions([]string{string(entities.PermissionManageMembers)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	// Filter by root
	rows, err := repo.ListRoles(ctx, &rootID)
	if err != nil {
		t.Fatalf("ListRoles: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(rows))
	}
	// Ordered by display_name ascending
	if rows[0].DisplayName != "Alpha" || rows[1].DisplayName != "Zeta" {
		t.Fatalf("expected ordered [Alpha, Zeta], got [%s, %s]", rows[0].DisplayName, rows[1].DisplayName)
	}

	// No filter => 3
	all, err := repo.ListRoles(ctx, nil)
	if err != nil {
		t.Fatalf("ListRoles(all): %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 roles, got %d", len(all))
	}
}

func TestCreateGetUpdateRole(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	created, err := repo.CreateRole(ctx, rootID, "lead", "Lead", []entities.Permission{entities.PermissionCreateProject})
	if err != nil {
		t.Fatalf("CreateRole: %v", err)
	}

	got, err := repo.GetRole(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRole: %v", err)
	}
	if got.Key != "lead" || got.DisplayName != "Lead" {
		t.Fatalf("unexpected role: %+v", got)
	}

	updated, err := repo.UpdateRole(ctx, created.ID, "Team lead", []entities.Permission{entities.PermissionManageDetails})
	if err != nil {
		t.Fatalf("UpdateRole: %v", err)
	}
	if updated.DisplayName != "Team lead" {
		t.Fatalf("expected display name updated, got %q", updated.DisplayName)
	}

	// Ensure persisted
	row, err := cli.OrganisationRole.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("ent get: %v", err)
	}
	if row.DisplayName != "Team lead" {
		t.Fatalf("expected persisted display name, got %q", row.DisplayName)
	}
	if !lo.Contains(row.Permissions, string(entities.PermissionManageDetails)) {
		t.Fatalf("expected manage_details permission persisted, got %v", row.Permissions)
	}
}

func TestDeleteRole_BlockedWhenScopeExists(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	// role
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("admin").
		SetDisplayName("Admin").
		SetPermissions([]string{string(entities.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	// scope referencing role
	_, err = cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}

	err = repo.DeleteRole(ctx, r.ID)
	if err == nil {
		t.Fatalf("expected error when deleting role in use")
	}

	// Ensure still exists
	exists, err := cli.OrganisationRole.Query().Where(entorgrole.IDEQ(r.ID)).Exist(ctx)
	if err != nil {
		t.Fatalf("exist query: %v", err)
	}
	if !exists {
		t.Fatalf("expected role to still exist")
	}
}

func TestCreateScope_IdempotentAndRequiresRoleForOrg(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	// Create role with key on root
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("contributor").
		SetDisplayName("Contributor").
		SetPermissions([]string{string(entities.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	s1, err := repo.CreateScope(ctx, "contributor", rootID)
	if err != nil {
		t.Fatalf("CreateScope: %v", err)
	}
	if s1.RoleID != r.ID || s1.RootNodeID != rootID {
		t.Fatalf("unexpected scope: %+v", s1)
	}

	// Second call should return existing scope (same id)
	s2, err := repo.CreateScope(ctx, "contributor", rootID)
	if err != nil {
		t.Fatalf("CreateScope again: %v", err)
	}
	if s1.ID != s2.ID {
		t.Fatalf("expected idempotent CreateScope to return same scope id")
	}

	// Wrong key for org => error
	_, err = repo.CreateScope(ctx, "missing_key", rootID)
	if err == nil {
		t.Fatalf("expected error for missing role key")
	}
}

func TestAddMembership_DeduplicatesAndRemoveMembership(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// Role + scope
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("member").
		SetDisplayName("Member").
		SetPermissions([]string{string(entities.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	sc, err := cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}

	m1, err := repo.AddMembership(ctx, personID, sc.ID)
	if err != nil {
		t.Fatalf("AddMembership: %v", err)
	}

	// Duplicate should error
	_, err = repo.AddMembership(ctx, personID, sc.ID)
	if err == nil {
		t.Fatalf("expected duplicate membership error")
	}

	// Remove and ensure gone
	if err := repo.RemoveMembership(ctx, m1.ID); err != nil {
		t.Fatalf("RemoveMembership: %v", err)
	}
	exists, err := cli.Membership.Query().Exist(ctx)
	if err != nil {
		t.Fatalf("exist query: %v", err)
	}
	if exists {
		t.Fatalf("expected no memberships after removal")
	}
}

func TestCreateScope_UsesRoleKeyScopedToOrg(t *testing.T) {
	cli := NewRoleClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := role.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)

	// Same key in two orgs is allowed; CreateScope must pick correct orgID.
	rRoot, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("same").
		SetDisplayName("SameRoot").
		SetPermissions([]string{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role root: %v", err)
	}
	rChild, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(childID).
		SetKey("same").
		SetDisplayName("SameChild").
		SetPermissions([]string{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role child: %v", err)
	}

	sRoot, err := repo.CreateScope(ctx, "same", rootID)
	if err != nil {
		t.Fatalf("CreateScope root: %v", err)
	}
	if sRoot.RoleID != rRoot.ID {
		t.Fatalf("expected scope to use root role id %s, got %s", rRoot.ID, sRoot.RoleID)
	}

	sChild, err := repo.CreateScope(ctx, "same", childID)
	if err != nil {
		t.Fatalf("CreateScope child: %v", err)
	}
	if sChild.RoleID != rChild.ID {
		t.Fatalf("expected scope to use child role id %s, got %s", rChild.ID, sChild.RoleID)
	}

	// Sanity that scopes persisted for each org
	c1, err := cli.RoleScope.Query().Where(entrolescope.RootNodeIDEQ(rootID)).Count(ctx)
	if err != nil {
		t.Fatalf("count scopes root: %v", err)
	}
	c2, err := cli.RoleScope.Query().Where(entrolescope.RootNodeIDEQ(childID)).Count(ctx)
	if err != nil {
		t.Fatalf("count scopes child: %v", err)
	}
	if c1 != 1 || c2 != 1 {
		t.Fatalf("expected 1 scope each, got root=%d child=%d", c1, c2)
	}
}
