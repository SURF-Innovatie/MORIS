package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	EnsureDefaultRoles(ctx context.Context) error
	ListRoles(ctx context.Context) ([]entities.OrganisationRole, error)

	CreateScope(ctx context.Context, roleKey string, rootNodeID uuid.UUID) (*entities.RoleScope, error)
	GetScope(ctx context.Context, id uuid.UUID) (*entities.RoleScope, error)
	AddMembership(ctx context.Context, personID uuid.UUID, roleScopeID uuid.UUID) (*entities.Membership, error)
	RemoveMembership(ctx context.Context, membershipID uuid.UUID) error

	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*entities.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) EnsureDefaultRoles(ctx context.Context) error {
	return s.repo.EnsureDefaultRoles(ctx)
}
func (s *service) ListRoles(ctx context.Context) ([]entities.OrganisationRole, error) {
	return s.repo.ListRoles(ctx)
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
