package errorlog

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/errorlog"
)

type Repository interface {
	Create(ctx context.Context, in errorlog.ErrorLogCreateInput) error
}
