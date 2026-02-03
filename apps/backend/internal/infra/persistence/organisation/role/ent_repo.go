package role

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entmembership "github.com/SURF-Innovatie/MORIS/ent/membership"
	entorgrole "github.com/SURF-Innovatie/MORIS/ent/organisationrole"
	entrolescope "github.com/SURF-Innovatie/MORIS/ent/rolescope"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
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

func (r *EntRepo) ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*rbac.OrganisationRole, error) {
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
	return transform.ToEntitiesPtr[rbac.OrganisationRole](rows), nil
}

func (r *EntRepo) CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error) {
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
	return transform.ToEntityPtr[rbac.OrganisationRole](row), nil
}

func (r *EntRepo) GetRole(ctx context.Context, roleID uuid.UUID) (*rbac.OrganisationRole, error) {
	row, err := r.cli.OrganisationRole.Get(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[rbac.OrganisationRole](row), nil
}

func (r *EntRepo) UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error) {
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
	return transform.ToEntityPtr[rbac.OrganisationRole](row), nil
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

func (r *EntRepo) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*rbac.RoleScope, error) {
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
		return transform.ToEntityPtr[rbac.RoleScope](existing), nil
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

	return transform.ToEntityPtr[rbac.RoleScope](row), nil
}

func (r *EntRepo) GetScope(ctx context.Context, id uuid.UUID) (*rbac.RoleScope, error) {
	row, err := r.cli.RoleScope.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[rbac.RoleScope](row), nil
}

func (r *EntRepo) AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*rbac.Membership, error) {
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

	return transform.ToEntityPtr[rbac.Membership](row), nil
}

func (r *EntRepo) GetMembership(ctx context.Context, membershipID uuid.UUID) (*rbac.Membership, error) {
	row, err := r.cli.Membership.Get(ctx, membershipID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[rbac.Membership](row), nil
}

func (r *EntRepo) RemoveMembership(ctx context.Context, membershipID uuid.UUID) error {
	return r.cli.Membership.DeleteOneID(membershipID).Exec(ctx)
}
