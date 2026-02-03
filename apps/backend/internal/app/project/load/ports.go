package load

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

type Cache interface {
	GetProject(ctx context.Context, id uuid.UUID) (*project.Project, error)
	SetProject(ctx context.Context, proj *project.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}
