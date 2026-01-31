package load

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Cache interface {
	GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	SetProject(ctx context.Context, proj *entities.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}
