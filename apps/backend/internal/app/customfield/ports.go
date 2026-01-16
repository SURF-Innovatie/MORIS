package customfield

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in CreateDefinitionInput) (*entities.CustomFieldDefinition, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsInOrg(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (bool, error)
	ListAvailableForNode(ctx context.Context, orgID uuid.UUID, category *entities.CustomFieldCategory) ([]entities.CustomFieldDefinition, error)
}
