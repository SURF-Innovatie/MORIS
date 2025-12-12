package organisation

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entmembership "github.com/SURF-Innovatie/MORIS/ent/membership"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type RBACService interface {
	EnsureDefaultRoles(ctx context.Context) error
	ListRoles(ctx context.Context) ([]entities.OrganisationRole, error)

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
}

type EffectiveMembership struct {
	MembershipID uuid.UUID
	PersonID     uuid.UUID

	RoleScopeID uuid.UUID
	ScopeRootID uuid.UUID

	RoleID         uuid.UUID
	RoleKey        string
	HasAdminRights bool
}

type rbacService struct {
	cli *ent.Client
}

func NewRBACService(cli *ent.Client) RBACService {
	return &rbacService{cli: cli}
}

func (s *rbacService) EnsureDefaultRoles(ctx context.Context) error {
	type def struct {
		key   string
		admin bool
	}

	defs := []def{
		{key: "admin", admin: true},
		{key: "researcher", admin: false},
		{key: "students", admin: false},
	}

	for _, d := range defs {
		existing, err := s.cli.OrganisationRole.
			Query().
			Where(entorgrole.KeyEQ(d.key)).
			Only(ctx)

		if err == nil {
			// update if needed
			if existing.HasAdminRights != d.admin {
				if _, err := s.cli.OrganisationRole.
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

		// create
		if _, err := s.cli.OrganisationRole.
			Create().
			SetKey(d.key).
			SetHasAdminRights(d.admin).
			Save(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *rbacService) ListRoles(ctx context.Context) ([]entities.OrganisationRole, error) {
	rows, err := s.cli.OrganisationRole.
		Query().
		Order(ent.Asc(entorgrole.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationRole, 0, len(rows))
	for _, r := range rows {
		out = append(out, entities.OrganisationRole{
			ID:             r.ID,
			Key:            r.Key,
			HasAdminRights: r.HasAdminRights,
		})
	}
	return out, nil
}

func (s *rbacService) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error) {
	role, err := s.cli.OrganisationRole.
		Query().
		Where(entorgrole.KeyEQ(roleKey)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// ensure node exists
	if _, err := s.cli.OrganisationNode.Get(ctx, rootNodeID); err != nil {
		return nil, err
	}

	row, err := s.cli.RoleScope.
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

func (s *rbacService) AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error) {
	if _, err := s.cli.Person.Get(ctx, personID); err != nil {
		return nil, err
	}
	if _, err := s.cli.RoleScope.Get(ctx, roleScopeID); err != nil {
		return nil, err
	}

	row, err := s.cli.Membership.
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

func (s *rbacService) RemoveMembership(ctx context.Context, membershipID uuid.UUID) error {
	return s.cli.Membership.DeleteOneID(membershipID).Exec(ctx)
}

func (s *rbacService) ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error) {
	// all ancestors of nodeID (including itself)
	ancestors, err := s.cli.OrganisationNodeClosure.
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

	scopes, err := s.cli.RoleScope.
		Query().
		Where(entrolescope.RootNodeIDIn(ancestorIDs...)).
		WithRole().
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(scopes) == 0 {
		return []EffectiveMembership{}, nil
	}

	scopeIDs := make([]uuid.UUID, 0, len(scopes))
	scopeByID := map[uuid.UUID]*ent.RoleScope{}
	for _, sc := range scopes {
		scopeIDs = append(scopeIDs, sc.ID)
		scopeByID[sc.ID] = sc
	}

	memberships, err := s.cli.Membership.
		Query().
		Where(entmembership.RoleScopeIDIn(scopeIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]EffectiveMembership, 0, len(memberships))
	for _, m := range memberships {
		sc := scopeByID[m.RoleScopeID]
		if sc == nil || sc.Edges.Role == nil {
			continue
		}
		out = append(out, EffectiveMembership{
			MembershipID:   m.ID,
			PersonID:       m.PersonID,
			RoleScopeID:    m.RoleScopeID,
			ScopeRootID:    sc.RootNodeID,
			RoleID:         sc.RoleID,
			RoleKey:        sc.Edges.Role.Key,
			HasAdminRights: sc.Edges.Role.HasAdminRights,
		})
	}

	return out, nil
}

func (s *rbacService) GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error) {
	rows, err := s.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		Order(ent.Asc(entclosure.FieldDepth)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		ancestorID := row.AncestorID

		adminScopes, err := s.cli.RoleScope.
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

		count, err := s.cli.Membership.
			Query().
			Where(entmembership.RoleScopeIDIn(scopeIDs...)).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			continue
		}

		n, err := s.cli.OrganisationNode.Get(ctx, ancestorID)
		if err != nil {
			return nil, err
		}
		return entities.OrganisationNodeFromEnt(n), nil
	}

	return nil, fmt.Errorf("no approval node found: ensure an admin membership exists in some ancestor scope")
}
