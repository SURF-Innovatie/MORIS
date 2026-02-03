package portfolio

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/portfolio"
	"github.com/google/uuid"
)

type Repository interface {
	GetByPersonID(ctx context.Context, personID uuid.UUID) (*portfolio.Portfolio, error)
	Upsert(ctx context.Context, portfolio portfolio.Portfolio) (*portfolio.Portfolio, error)
}
