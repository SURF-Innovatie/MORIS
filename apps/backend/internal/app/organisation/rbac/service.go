package organisation_rbac

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/google/uuid"
)

type Service interface {
	ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error)
	ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error)
	GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]rbac.Permission, error)

	GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*organisation.OrganisationNode, error)
	HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission rbac.Permission) (bool, error)
}

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{repo: repo}
}

func (s *service) ListEffectiveMemberships(ctx context.Context, nodeID uuid.UUID) ([]EffectiveMembership, error) {
	return s.repo.ListEffectiveMemberships(ctx, nodeID)
}
func (s *service) ListMyMemberships(ctx context.Context, personID uuid.UUID) ([]EffectiveMembership, error) {
	return s.repo.ListMyMemberships(ctx, personID)
}

func (s *service) GetApprovalNode(ctx context.Context, nodeID uuid.UUID) (*organisation.OrganisationNode, error) {
	return s.repo.GetApprovalNode(ctx, nodeID)
}
func (s *service) HasAdminAccess(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	return s.repo.HasAdminAccess(ctx, personID, nodeID)
}
func (s *service) HasPermission(ctx context.Context, personID uuid.UUID, nodeID uuid.UUID, permission rbac.Permission) (bool, error) {
	return s.repo.HasPermission(ctx, personID, nodeID, permission)
}

func (s *service) GetMyPermissions(ctx context.Context, userID, nodeID uuid.UUID) ([]rbac.Permission, error) {
	return s.repo.GetMyPermissions(ctx, userID, nodeID)
}
