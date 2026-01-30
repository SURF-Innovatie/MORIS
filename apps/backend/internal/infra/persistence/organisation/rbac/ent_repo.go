package rbac

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entmembership "github.com/SURF-Innovatie/MORIS/ent/membership"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	organisation_rbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

// EnsureDefaultRoles upserts/ensures keys exist + admin flags are correct.
// EnsureDefaultRoles is deprecated as roles are now per-organisation.
func (r *EntRepo) EnsureDefaultRoles(ctx context.Context) error {
	return nil
}

func (r *EntRepo) ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*entities.OrganisationRole, error) {
	q := r.cli.OrganisationRole.Query()
	if orgID != nil {
		q = q.Where(entorgrole.OrganisationNodeIDEQ(*orgID))
	}
	rows, err := q.
		Order(ent.Asc(entorgrole.FieldDisplayName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[entities.OrganisationRole](rows), nil
}

func (r *EntRepo) GetMyPermissions(ctx context.Context, personID, nodeID uuid.UUID) ([]role.Permission, error) {
	// Check if user is sysadmin
	user, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).First(ctx)
	if err == nil && user.IsSysAdmin {
		return role.AllPermissions, nil
	}

	// all ancestors of nodeID (including itself)
	ancestors, err := r.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	ancestorIDs := lo.Map(ancestors, func(a *ent.OrganisationNodeClosure, _ int) uuid.UUID {
		return a.AncestorID
	})

	// scopes defined on ancestors
	scopes, err := r.cli.RoleScope.
		Query().
		Where(entrolescope.RootNodeIDIn(ancestorIDs...)).
		WithRole().
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(scopes) == 0 {
		return []role.Permission{}, nil
	}

	scopeIDs := lo.Map(scopes, func(sc *ent.RoleScope, _ int) uuid.UUID { return sc.ID })
	scopeByID := lo.KeyBy(scopes, func(sc *ent.RoleScope) uuid.UUID { return sc.ID })

	memberships, err := r.cli.Membership.
		Query().
		Where(
			entmembership.PersonIDEQ(personID),
			entmembership.RoleScopeIDIn(scopeIDs...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	permissions := make(map[role.Permission]struct{})
	for _, m := range memberships {
		scope, ok := scopeByID[m.RoleScopeID]
		if !ok {
			continue // Should not happen
		}
		if scope.Edges.Role == nil {
			continue // Should not happen
		}
		for _, p := range scope.Edges.Role.Permissions {
			permissions[role.Permission(p)] = struct{}{}
		}
	}

	result := lo.Keys(permissions)
	return result, nil
}

func (r *EntRepo) CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error) {
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}

	row, err := r.cli.OrganisationRole.
		Create().
		SetOrganisationNodeID(orgID).
		SetKey(key).
		SetDisplayName(displayName).
		SetPermissions(perms).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.OrganisationRole](row), nil
}

func (r *EntRepo) GetRole(ctx context.Context, roleID uuid.UUID) (*entities.OrganisationRole, error) {
	row, err := r.cli.OrganisationRole.Get(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.OrganisationRole](row), nil
}

func (r *EntRepo) UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error) {
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}

	row, err := r.cli.OrganisationRole.
		UpdateOneID(roleID).
		SetDisplayName(displayName).
		SetPermissions(perms).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.OrganisationRole](row), nil
}

func (r *EntRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	// Check if used in scopes
	inUse, err := r.cli.RoleScope.Query().Where(entrolescope.RoleIDEQ(roleID)).Exist(ctx)
	if err != nil {
		return err
	}
	if inUse {
		return fmt.Errorf("role is in use by memberships and cannot be deleted")
	}

	return r.cli.OrganisationRole.DeleteOneID(roleID).Exec(ctx)
}

func (r *EntRepo) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error) {
	// Role Key is not unique globally anymore.
	// We assume CreateScope is called with a key that is valid for the given rootNodeID (org).
	role, err := r.cli.OrganisationRole.
		Query().
		Where(
			entorgrole.KeyEQ(roleKey),
			entorgrole.OrganisationNodeIDEQ(rootNodeID),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("role %q not found for org %s: %w", roleKey, rootNodeID, err)
	}

	if _, err := r.cli.OrganisationNode.Get(ctx, rootNodeID); err != nil {
		return nil, err
	}

	existing, err := r.cli.RoleScope.
		Query().
		Where(
			entrolescope.RoleIDEQ(role.ID),
			entrolescope.RootNodeIDEQ(rootNodeID),
		).
		Only(ctx)

	if err == nil {
		return transform.ToEntityPtr[entities.RoleScope](existing), nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}

	row, err := r.cli.RoleScope.
		Create().
		SetRoleID(role.ID).
		SetRootNodeID(rootNodeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.RoleScope](row), nil
}

func (r *EntRepo) GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error) {
	row, err := r.cli.RoleScope.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.RoleScope](row), nil
}

func (r *EntRepo) AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error) {
	if _, err := r.cli.Person.Get(ctx, personID); err != nil {
		return nil, err
	}
	if _, err := r.cli.RoleScope.Get(ctx, roleScopeID); err != nil {
		return nil, err
	}

	exists, err := r.cli.Membership.
		Query().
		Where(
			entmembership.PersonIDEQ(personID),
			entmembership.RoleScopeIDEQ(roleScopeID),
		).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("membership already exists")
	}

	row, err := r.cli.Membership.
		Create().
		SetPersonID(personID).
		SetRoleScopeID(roleScopeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.Membership](row), nil
}

func (r *EntRepo) GetMembership(ctx context.Context, membershipID uuid.UUID) (*entities.Membership, error) {
	row, err := r.cli.Membership.Get(ctx, membershipID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Membership](row), nil
}

func (r *EntRepo) RemoveMembership(ctx context.Context, membershipID uuid.UUID) error {
	return r.cli.Membership.DeleteOneID(membershipID).Exec(ctx)
}

func (r *EntRepo) ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]organisation_rbac.EffectiveMembership, error) {
	// all ancestors of nodeID (including itself)
	ancestors, err := r.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	ancestorIDs := make([]uuid.UUID, 0, len(ancestors))
	for _, a := range ancestors {
		ancestorIDs = append(ancestorIDs, a.AncestorID)
	}

	// scopes defined on ancestors
	scopes, err := r.cli.RoleScope.
		Query().
		Where(entrolescope.RootNodeIDIn(ancestorIDs...)).
		WithRole().
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(scopes) == 0 {
		return []organisation_rbac.EffectiveMembership{}, nil
	}

	scopeIDs := make([]uuid.UUID, 0, len(scopes))
	scopeByID := map[uuid.UUID]*ent.RoleScope{}
	for _, sc := range scopes {
		scopeIDs = append(scopeIDs, sc.ID)
		scopeByID[sc.ID] = sc
	}

	memberships, err := r.cli.Membership.
		Query().
		Where(entmembership.RoleScopeIDIn(scopeIDs...)).
		WithPerson().
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]organisation_rbac.EffectiveMembership, 0, len(memberships))
	for _, m := range memberships {
		sc := scopeByID[m.RoleScopeID]
		if sc == nil || sc.Edges.Role == nil {
			continue
		}

		n, err := r.cli.OrganisationNode.Get(ctx, sc.RootNodeID)
		if err != nil {
			return nil, err
		}
		scopeRoot := &entities.OrganisationNode{
			ID:       n.ID,
			ParentID: n.ParentID,
			Name:     n.Name,
		}

		p := m.Edges.Person
		if p == nil {
			continue
		}

		out = append(out, organisation_rbac.EffectiveMembership{
			MembershipID:          m.ID,
			PersonID:              m.PersonID,
			RoleScopeID:           m.RoleScopeID,
			ScopeRootOrganisation: scopeRoot,
			RoleID:                sc.RoleID,
			RoleKey:               sc.Edges.Role.Key,
			Permissions:           toPermissions(sc.Edges.Role.Permissions),
			CustomFields:          getCustomFieldsForNode(p.OrgCustomFields, nodeID),
			Person:                transform.ToEntity[entities.Person](p),
		})
	}

	return out, nil
}

func (r *EntRepo) ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]organisation_rbac.EffectiveMembership, error) {
	memberships, err := r.cli.Membership.
		Query().
		Where(entmembership.PersonIDEQ(personID)).
		WithRoleScope(func(q *ent.RoleScopeQuery) { q.WithRole() }).
		WithPerson().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Check if user is sysadmin to override rights
	user, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).First(ctx)
	isSysAdmin := false
	if err == nil && user != nil {
		isSysAdmin = user.IsSysAdmin
	}

	out := make([]organisation_rbac.EffectiveMembership, 0, len(memberships))
	for _, m := range memberships {
		sc := m.Edges.RoleScope
		if sc == nil || sc.Edges.Role == nil {
			continue
		}
		p := m.Edges.Person
		if p == nil {
			continue
		}
		n, err := r.cli.OrganisationNode.Get(ctx, sc.RootNodeID)
		if err != nil {
			return nil, err
		}
		scopeRoot := transform.ToEntityPtr[entities.OrganisationNode](n)

		out = append(out, organisation_rbac.EffectiveMembership{
			MembershipID:          m.ID,
			PersonID:              m.PersonID,
			RoleScopeID:           m.RoleScopeID,
			ScopeRootOrganisation: scopeRoot,
			RoleID:                sc.RoleID,
			RoleKey:               sc.Edges.Role.Key,
			Permissions:           toPermissions(sc.Edges.Role.Permissions),
			HasAdminRights:        isSysAdmin || hasAdminPermission(sc.Edges.Role.Permissions),
			CustomFields:          getCustomFieldsForNode(p.OrgCustomFields, sc.RootNodeID),
		})
	}
	return out, nil
}

func (r *EntRepo) GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error) {
	rows, err := r.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		Order(ent.Asc(entclosure.FieldDepth)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		ancestorID := row.AncestorID

		scopes, err := r.cli.RoleScope.
			Query().
			Where(entrolescope.RootNodeIDEQ(ancestorID)).
			WithRole().
			All(ctx)
		if err != nil {
			return nil, err
		}

		// Filter scopes for Admin role (manage_details permission)
		var adminScopeIDs []uuid.UUID
		for _, sc := range scopes { // Note: iterate scopes
			if sc.Edges.Role == nil {
				continue
			}
			for _, p := range sc.Edges.Role.Permissions {
				if role.Permission(p) == role.PermissionManageDetails {
					adminScopeIDs = append(adminScopeIDs, sc.ID)
					break
				}
			}
		}

		if len(adminScopeIDs) == 0 {
			continue
		}

		count, err := r.cli.Membership.
			Query().
			Where(entmembership.RoleScopeIDIn(adminScopeIDs...)).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			continue
		}

		n, err := r.cli.OrganisationNode.Get(ctx, ancestorID)
		if err != nil {
			return nil, err
		}
		return transform.ToEntityPtr[entities.OrganisationNode](n), nil
	}

	return nil, fmt.Errorf("no approval node found: ensure an admin membership exists in some ancestor scope")
}

func (r *EntRepo) HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	return r.HasPermission(ctx, personID, nodeID, role.PermissionManageDetails)
}

func (r *EntRepo) HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission role.Permission) (bool, error) {
	user, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).First(ctx)
	if err != nil {
		return false, err
	}
	if user.IsSysAdmin {
		return true, nil
	}

	ancestors, err := r.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		All(ctx)
	if err != nil {
		return false, err
	}

	ancestorIDs := make([]uuid.UUID, 0, len(ancestors))
	for _, a := range ancestors {
		ancestorIDs = append(ancestorIDs, a.AncestorID)
	}

	scopes, err := r.cli.RoleScope.
		Query().
		Where(entrolescope.RootNodeIDIn(ancestorIDs...)).
		WithRole().
		All(ctx)
	if err != nil {
		return false, err
	}

	var validScopeIDs []uuid.UUID
	for _, sc := range scopes {
		if sc.Edges.Role == nil {
			continue
		}
		for _, p := range sc.Edges.Role.Permissions {
			if role.Permission(p) == permission {
				validScopeIDs = append(validScopeIDs, sc.ID)
				break
			}
		}
	}

	if len(validScopeIDs) == 0 {
		return false, nil
	}

	count, err := r.cli.Membership.
		Query().
		Where(
			entmembership.RoleScopeIDIn(validScopeIDs...),
			entmembership.PersonIDEQ(personID),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func toPermissions(s []string) []role.Permission {
	return lo.Map(s, func(v string, _ int) role.Permission {
		return role.Permission(v)
	})
}

// Helper to safely get custom fields for a node
func getCustomFieldsForNode(fields map[string]interface{}, nodeID uuid.UUID) map[string]interface{} {
	if fields == nil {
		return nil
	}
	if v, ok := fields[nodeID.String()]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}
func hasAdminPermission(perms []string) bool {
	return lo.SomeBy(perms, func(p string) bool {
		return p == string(role.PermissionManageDetails) || p == string(role.PermissionManageMembers)
	})
}
