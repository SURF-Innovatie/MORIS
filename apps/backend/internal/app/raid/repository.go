package raid

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, raidID string) (*entities.RAiDInfo, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.RAiDInfo, error)
	Create(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error)
	Update(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error)
	Delete(ctx context.Context, raidID string) error
}
