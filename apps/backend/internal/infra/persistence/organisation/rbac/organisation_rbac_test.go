package rbac_test

import (
	"context"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/enttest"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/rbac"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

func newRBACClient(t *testing.T) *ent.Client {
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
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)

	rootID, childID := seedOrgTree(t, cli)

	// Create two roles in root (unordered display names)
	_, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("b").
		SetDisplayName("Zeta").
		SetPermissions([]string{string(role.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	_, err = cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("a").
		SetDisplayName("Alpha").
		SetPermissions([]string{string(role.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	// Another role in child org
	_, err = cli.OrganisationRole.Create().
		SetOrganisationNodeID(childID).
		SetKey("c").
		SetDisplayName("Other").
		SetPermissions([]string{string(role.PermissionManageMembers)}).
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
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	created, err := repo.CreateRole(ctx, rootID, "lead", "Lead", []role.Permission{role.PermissionCreateProject})
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

	updated, err := repo.UpdateRole(ctx, created.ID, "Team lead", []role.Permission{role.PermissionManageDetails})
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
	if !lo.Contains(row.Permissions, string(role.PermissionManageDetails)) {
		t.Fatalf("expected manage_details permission persisted, got %v", row.Permissions)
	}
}

func TestDeleteRole_BlockedWhenScopeExists(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	// role
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("admin").
		SetDisplayName("Admin").
		SetPermissions([]string{string(role.PermissionManageDetails)}).
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
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)

	// Create role with key on root
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("contributor").
		SetDisplayName("Contributor").
		SetPermissions([]string{string(role.PermissionCreateProject)}).
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
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// Role + scope
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("member").
		SetDisplayName("Member").
		SetPermissions([]string{string(role.PermissionCreateProject)}).
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

func TestHasPermission_TrueViaAncestorScope_FalseOtherwise(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// role on ROOT with manage_details
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("admin").
		SetDisplayName("Admin").
		SetPermissions([]string{string(role.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	// scope rooted at ROOT
	sc, err := cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}
	// membership for person in that scope
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(sc.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership: %v", err)
	}

	ok, err := repo.HasPermission(ctx, personID, childID, role.PermissionManageDetails)
	if err != nil {
		t.Fatalf("HasPermission: %v", err)
	}
	if !ok {
		t.Fatalf("expected permission via ancestor scope")
	}

	ok, err = repo.HasPermission(ctx, personID, childID, role.PermissionCreateProject)
	if err != nil {
		t.Fatalf("HasPermission: %v", err)
	}
	if ok {
		t.Fatalf("expected no create_project permission")
	}
}

func TestHasPermission_SysAdminAlwaysTrue(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	_, childID := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, true)

	ok, err := repo.HasPermission(ctx, personID, childID, role.PermissionManageDetails)
	if err != nil {
		t.Fatalf("HasPermission: %v", err)
	}
	if !ok {
		t.Fatalf("expected sysadmin to have permission")
	}

	perms, err := repo.GetMyPermissions(ctx, personID, childID)
	if err != nil {
		t.Fatalf("GetMyPermissions: %v", err)
	}
	if len(perms) != len(role.AllPermissions) {
		t.Fatalf("expected all permissions for sysadmin, got %d", len(perms))
	}
}

func TestGetMyPermissions_CollectsUnionFromMemberships(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// Two roles on root with different permissions
	r1, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("r1").
		SetDisplayName("R1").
		SetPermissions([]string{string(role.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role1: %v", err)
	}
	r2, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("r2").
		SetDisplayName("R2").
		SetPermissions([]string{string(role.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role2: %v", err)
	}

	// scopes rooted at root
	s1, err := cli.RoleScope.Create().SetRoleID(r1.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope1: %v", err)
	}
	s2, err := cli.RoleScope.Create().SetRoleID(r2.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope2: %v", err)
	}

	// memberships for person in both scopes
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(s1.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership1: %v", err)
	}
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(s2.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership2: %v", err)
	}

	perms, err := repo.GetMyPermissions(ctx, personID, childID)
	if err != nil {
		t.Fatalf("GetMyPermissions: %v", err)
	}

	if !lo.Contains(perms, role.PermissionManageDetails) || !lo.Contains(perms, role.PermissionCreateProject) {
		t.Fatalf("expected union of permissions, got %v", perms)
	}

	// Ensure it does NOT depend on any specific order
	if len(lo.Uniq(perms)) != len(perms) {
		t.Fatalf("expected unique permissions, got %v", perms)
	}
}

func TestListEffectiveMemberships_ReturnsMembershipsForAncestorScopes(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// role on root, scope on root, membership
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("member").
		SetDisplayName("Member").
		SetPermissions([]string{string(role.PermissionCreateProject)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	sc, err := cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(sc.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership: %v", err)
	}

	eff, err := repo.ListEffectiveMemberships(ctx, childID)
	if err != nil {
		t.Fatalf("ListEffectiveMemberships: %v", err)
	}
	if len(eff) != 1 {
		t.Fatalf("expected 1 effective membership, got %d", len(eff))
	}

	if eff[0].RoleKey != "member" {
		t.Fatalf("expected rolekey 'member', got %q", eff[0].RoleKey)
	}
	if eff[0].ScopeRootOrganisation == nil || eff[0].ScopeRootOrganisation.ID != rootID {
		t.Fatalf("expected scope root to be root org")
	}
}

func TestListMyMemberships_SetsHasAdminRightsForSysAdmin(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, _ := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, true)

	// role + scope + membership (role perms irrelevant for sysadmin flag)
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("x").
		SetDisplayName("X").
		SetPermissions([]string{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	sc, err := cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(sc.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership: %v", err)
	}

	out, err := repo.ListMyMemberships(ctx, personID)
	if err != nil {
		t.Fatalf("ListMyMemberships: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 membership, got %d", len(out))
	}
	if !out[0].HasAdminRights {
		t.Fatalf("expected HasAdminRights=true for sysadmin")
	}

	// sanity: repo checks sysadmin via user table
	exists, err := cli.User.Query().Where(entuser.PersonIDEQ(personID)).Exist(ctx)
	if err != nil {
		t.Fatalf("user exist query: %v", err)
	}
	if !exists {
		t.Fatalf("expected user row to exist")
	}
}

func TestGetApprovalNode_FindsNearestAncestorWithAdminMembership(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	rootID, childID := seedOrgTree(t, cli)
	personID, _ := seedPersonUser(t, cli, false)

	// Create admin role on root with manage_details
	r, err := cli.OrganisationRole.Create().
		SetOrganisationNodeID(rootID).
		SetKey("admin").
		SetDisplayName("Admin").
		SetPermissions([]string{string(role.PermissionManageDetails)}).
		Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	sc, err := cli.RoleScope.Create().SetRoleID(r.ID).SetRootNodeID(rootID).Save(ctx)
	if err != nil {
		t.Fatalf("create scope: %v", err)
	}
	_, err = cli.Membership.Create().SetPersonID(personID).SetRoleScopeID(sc.ID).Save(ctx)
	if err != nil {
		t.Fatalf("create membership: %v", err)
	}

	node, err := repo.GetApprovalNode(ctx, childID)
	if err != nil {
		t.Fatalf("GetApprovalNode: %v", err)
	}
	if node.ID != rootID {
		t.Fatalf("expected approval node to be root, got %s", node.ID)
	}

	// Also validate ordering by depth works by ensuring closure rows exist
	_, err = cli.OrganisationNodeClosure.Query().
		Where(entclosure.DescendantIDEQ(childID)).
		Order(ent.Asc(entclosure.FieldDepth)).
		All(ctx)
	if err != nil {
		t.Fatalf("closure query sanity: %v", err)
	}
}

func TestGetApprovalNode_ErrorsWhenNoAdminMembershipAnywhere(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
	_, childID := seedOrgTree(t, cli)

	_, err := repo.GetApprovalNode(ctx, childID)
	if err == nil {
		t.Fatalf("expected error when no approval node exists")
	}
}

func TestCreateScope_UsesRoleKeyScopedToOrg(t *testing.T) {
	cli := newRBACClient(t)
	defer cli.Close()
	ctx := context.Background()

	repo := rbac.NewEntRepo(cli)
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
