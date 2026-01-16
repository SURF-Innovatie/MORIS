package errorlog

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Repository interface {
	Create(ctx context.Context, in entities.ErrorLogCreateInput) error
}
