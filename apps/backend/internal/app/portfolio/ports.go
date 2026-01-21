package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	GetByPersonID(ctx context.Context, personID uuid.UUID) (*entities.Portfolio, error)
	Upsert(ctx context.Context, portfolio entities.Portfolio) (*entities.Portfolio, error)
}
