package rbac

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entmembership "github.com/SURF-Innovatie/MORIS/ent/membership"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	organisation_rbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) GetMyPermissions(ctx context.Context, personID, nodeID uuid.UUID) ([]rbac.Permission, error) {
	// Check if user is sysadmin
	user, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).First(ctx)
	if err == nil && user.IsSysAdmin {
		return rbac.AllPermissions, nil
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
		return []rbac.Permission{}, nil
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

	permissions := make(map[rbac.Permission]struct{})
	for _, m := range memberships {
		scope, ok := scopeByID[m.RoleScopeID]
		if !ok {
			continue // Should not happen
		}
		if scope.Edges.Role == nil {
			continue // Should not happen
		}
		for _, p := range scope.Edges.Role.Permissions {
			permissions[rbac.Permission(p)] = struct{}{}
		}
	}

	result := lo.Keys(permissions)
	return result, nil
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
		scopeRoot := &organisation.OrganisationNode{
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
			Person:                transform.ToEntity[identity.Person](p),
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
		scopeRoot := transform.ToEntityPtr[organisation.OrganisationNode](n)

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

func (r *EntRepo) GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*organisation.OrganisationNode, error) {
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
				if rbac.Permission(p) == rbac.PermissionManageDetails {
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
		return transform.ToEntityPtr[organisation.OrganisationNode](n), nil
	}

	return nil, fmt.Errorf("no approval node found: ensure an admin membership exists in some ancestor scope")
}

func (r *EntRepo) HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	return r.HasPermission(ctx, personID, nodeID, rbac.PermissionManageDetails)
}

func (r *EntRepo) HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission rbac.Permission) (bool, error) {
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
			if rbac.Permission(p) == permission {
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

func toPermissions(s []string) []rbac.Permission {
	return lo.Map(s, func(v string, _ int) rbac.Permission {
		return rbac.Permission(v)
	})
}

// Helper to safely get custom fields for a node
func getCustomFieldsForNode(fields map[string]any, nodeID uuid.UUID) map[string]any {
	if fields == nil {
		return nil
	}
	if v, ok := fields[nodeID.String()]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}
func hasAdminPermission(perms []string) bool {
	return lo.SomeBy(perms, func(p string) bool {
		return p == string(rbac.PermissionManageDetails) || p == string(rbac.PermissionManageMembers)
	})
}
