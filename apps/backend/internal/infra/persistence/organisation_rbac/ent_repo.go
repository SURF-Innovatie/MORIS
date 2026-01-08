package organisation_rbac

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entmembership "github.com/SURF-Innovatie/MORIS/ent/membership"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

// EnsureDefaultRoles upserts/ensures keys exist + admin flags are correct.
func (r *EntRepo) EnsureDefaultRoles(ctx context.Context) error {
	type def struct {
		key   string
		admin bool
	}

	defs := []def{
		{key: "admin", admin: true},
		{key: "researcher", admin: false},
		{key: "student", admin: false},
	}

	for _, d := range defs {
		existing, err := r.cli.OrganisationRole.
			Query().
			Where(entorgrole.KeyEQ(d.key)).
			Only(ctx)

		if err == nil {
			if existing.HasAdminRights != d.admin {
				if _, err := r.cli.OrganisationRole.
					UpdateOneID(existing.ID).
					SetHasAdminRights(d.admin).
					Save(ctx); err != nil {
					return err
				}
			}
			continue
		}
		if !ent.IsNotFound(err) {
			return err
		}

		if _, err := r.cli.OrganisationRole.
			Create().
			SetKey(d.key).
			SetHasAdminRights(d.admin).
			Save(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *EntRepo) ListRoles(ctx context.Context) ([]entities.OrganisationRole, error) {
	rows, err := r.cli.OrganisationRole.
		Query().
		Order(ent.Asc(entorgrole.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationRole, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.OrganisationRole{
			ID:             row.ID,
			Key:            row.Key,
			HasAdminRights: row.HasAdminRights,
		})
	}
	return out, nil
}

func (r *EntRepo) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error) {
	role, err := r.cli.OrganisationRole.
		Query().
		Where(entorgrole.KeyEQ(roleKey)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
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
		return &entities.RoleScope{
			ID:         existing.ID,
			RoleID:     existing.RoleID,
			RootNodeID: existing.RootNodeID,
		}, nil
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

	return &entities.RoleScope{
		ID:         row.ID,
		RoleID:     row.RoleID,
		RootNodeID: row.RootNodeID,
	}, nil
}

func (r *EntRepo) GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error) {
	row, err := r.cli.RoleScope.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entities.RoleScope{
		ID:         row.ID,
		RoleID:     row.RoleID,
		RootNodeID: row.RootNodeID,
	}, nil
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

	return &entities.Membership{
		ID:          row.ID,
		PersonID:    row.PersonID,
		RoleScopeID: row.RoleScopeID,
	}, nil
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
			HasAdminRights:        sc.Edges.Role.HasAdminRights,
			CustomFields:          getCustomFieldsForNode(p.OrgCustomFields, nodeID),
			Person: entities.Person{
				ID:          p.ID,
				Name:        p.Name,
				ORCiD:       &p.OrcidID,
				GivenName:   p.GivenName,
				FamilyName:  p.FamilyName,
				Email:       p.Email,
				AvatarUrl:   p.AvatarURL,
				Description: p.Description,
			},
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
		scopeRoot := &entities.OrganisationNode{
			ID:       n.ID,
			ParentID: n.ParentID,
			Name:     n.Name,
		}

		out = append(out, organisation_rbac.EffectiveMembership{
			MembershipID:          m.ID,
			PersonID:              m.PersonID,
			RoleScopeID:           m.RoleScopeID,
			ScopeRootOrganisation: scopeRoot,
			RoleID:                sc.RoleID,
			RoleKey:               sc.Edges.Role.Key,
			HasAdminRights:        sc.Edges.Role.HasAdminRights,
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

		adminScopes, err := r.cli.RoleScope.
			Query().
			Where(
				entrolescope.RootNodeIDEQ(ancestorID),
				entrolescope.HasRoleWith(entorgrole.HasAdminRightsEQ(true)),
			).
			All(ctx)
		if err != nil {
			return nil, err
		}
		if len(adminScopes) == 0 {
			continue
		}

		scopeIDs := make([]uuid.UUID, 0, len(adminScopes))
		for _, sc := range adminScopes {
			scopeIDs = append(scopeIDs, sc.ID)
		}

		count, err := r.cli.Membership.
			Query().
			Where(entmembership.RoleScopeIDIn(scopeIDs...)).
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
		return &entities.OrganisationNode{
			ID:       n.ID,
			ParentID: n.ParentID,
			Name:     n.Name,
		}, nil
	}

	return nil, fmt.Errorf("no approval node found: ensure an admin membership exists in some ancestor scope")
}

func (r *EntRepo) HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error) {
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

	adminScopes, err := r.cli.RoleScope.
		Query().
		Where(
			entrolescope.RootNodeIDIn(ancestorIDs...),
			entrolescope.HasRoleWith(entorgrole.HasAdminRightsEQ(true)),
		).
		All(ctx)
	if err != nil {
		return false, err
	}
	if len(adminScopes) == 0 {
		return false, nil
	}

	scopeIDs := make([]uuid.UUID, 0, len(adminScopes))
	for _, sc := range adminScopes {
		scopeIDs = append(scopeIDs, sc.ID)
	}

	count, err := r.cli.Membership.
		Query().
		Where(
			entmembership.RoleScopeIDIn(scopeIDs...),
			entmembership.PersonIDEQ(personID),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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


