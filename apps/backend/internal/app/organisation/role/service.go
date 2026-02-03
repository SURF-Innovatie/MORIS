package role

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
)

type Service interface {
	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*rbac.OrganisationRole, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*rbac.OrganisationRole, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*rbac.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*rbac.RoleScope, error)

	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*rbac.Membership, error)
	GetMembership(ctx context.Context, membershipID uuid.UUID) (*rbac.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	EnsureDefaultRoles(ctx context.Context) error
}

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{repo: repo}
}

func (s *service) EnsureDefaultRoles(ctx context.Context) error {
	return s.repo.EnsureDefaultRoles(ctx)
}

func (s *service) ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*rbac.OrganisationRole, error) {
	return s.repo.ListRoles(ctx, orgID)
}

func (s *service) CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error) {
	// TODO: Add logic to create RoleScope automatically if needed, or leave it to handler?
	role, err := s.repo.CreateRole(ctx, orgID, key, displayName, permissions)
	if err != nil {
		return nil, err
	}
	// Auto-create scope for this role on the same org
	_, err = s.repo.CreateScope(ctx, key, orgID)
	return role, nil
}

func (s *service) GetRole(ctx context.Context, roleID uuid.UUID) (*rbac.OrganisationRole, error) {
	return s.repo.GetRole(ctx, roleID)
}

func (s *service) UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []rbac.Permission) (*rbac.OrganisationRole, error) {
	return s.repo.UpdateRole(ctx, roleID, displayName, permissions)
}

func (s *service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	// TODO: Check if role has active memberships?
	// The plan said "Block if used".
	// Implementation should be in Repo or here.
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *service) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*rbac.RoleScope, error) {
	return s.repo.CreateScope(ctx, roleKey, rootNodeID)
}

func (s *service) GetScope(ctx context.Context, id uuid.UUID) (*rbac.RoleScope, error) {
	return s.repo.GetScope(ctx, id)
}

func (s *service) AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*rbac.Membership, error) {
	return s.repo.AddMembership(ctx, personID, roleScopeID)
}

func (s *service) GetMembership(ctx context.Context, membershipID uuid.UUID) (*rbac.Membership, error) {
	return s.repo.GetMembership(ctx, membershipID)
}
func (s *service) RemoveMembership(ctx context.Context, membershipID uuid.UUID) error {
	return s.repo.RemoveMembership(ctx, membershipID)
}
