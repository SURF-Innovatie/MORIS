package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	EnsureDefaultRoles(ctx context.Context) error
	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*entities.OrganisationRole, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*entities.OrganisationRole, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error)
	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	GetMembership(ctx context.Context, membershipID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]role.Permission, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission role.Permission) (bool, error)

	AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error)
	IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error)
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
func (s *service) ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*entities.OrganisationRole, error) {
	return s.repo.ListRoles(ctx, orgID)
}

func (s *service) CreateRole(ctx context.Context, orgID uuid.UUID, key, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error) {
	// TODO: Add logic to create RoleScope automatically if needed, or leave it to handler?
	role, err := s.repo.CreateRole(ctx, orgID, key, displayName, permissions)
	if err != nil {
		return nil, err
	}
	// Auto-create scope for this role on the same org
	_, err = s.repo.CreateScope(ctx, key, orgID)
	return role, nil
}

func (s *service) GetRole(ctx context.Context, roleID uuid.UUID) (*entities.OrganisationRole, error) {
	return s.repo.GetRole(ctx, roleID)
}

func (s *service) UpdateRole(ctx context.Context, roleID uuid.UUID, displayName string, permissions []role.Permission) (*entities.OrganisationRole, error) {
	return s.repo.UpdateRole(ctx, roleID, displayName, permissions)
}

func (s *service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	// TODO: Check if role has active memberships?
	// The plan said "Block if used".
	// Implementation should be in Repo or here.
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *service) CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error) {
	return s.repo.CreateScope(ctx, roleKey, rootNodeID)
}
func (s *service) GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error) {
	return s.repo.GetScope(ctx, id)
}

func (s *service) AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error) {
	return s.repo.AddMembership(ctx, personID, roleScopeID)
}
func (s *service) GetMembership(ctx context.Context, membershipID uuid.UUID) (*entities.Membership, error) {
	return s.repo.GetMembership(ctx, membershipID)
}
func (s *service) RemoveMembership(ctx context.Context, membershipID uuid.UUID) error {
	return s.repo.RemoveMembership(ctx, membershipID)
}

func (s *service) ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error) {
	return s.repo.ListEffectiveMemberships(ctx, nodeID)
}
func (s *service) ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error) {
	return s.repo.ListMyMemberships(ctx, personID)
}

func (s *service) GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error) {
	return s.repo.GetApprovalNode(ctx, nodeID)
}
func (s *service) HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	return s.repo.HasAdminAccess(ctx, personID, nodeID)
}
func (s *service) HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission role.Permission) (bool, error) {
	return s.repo.HasPermission(ctx, personID, nodeID, permission)
}

func (s *service) GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]role.Permission, error) {
	return s.repo.GetMyPermissions(ctx, userID, nodeID)
}

func (s *service) AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.AncestorIDs(ctx, nodeID)
}

func (s *service) IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error) {
	return s.repo.IsAncestor(ctx, ancestorID, descendantID)
}
