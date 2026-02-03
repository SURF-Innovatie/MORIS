package errorlog

import (
	"context"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/errorlog"
	"github.com/google/uuid"
)

type Service interface {
	Log(ctx context.Context, userID *uuid.UUID, method, route string, statusCode int, errMsg string, stackTrace *string)
}

type service struct {
	repo    Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	s := &service{
		repo:    repo,
		timeout: 5 * time.Second,
	}
	return s
}

func (s *service) Log(_ context.Context, userID *uuid.UUID, method, route string, statusCode int, errMsg string, stackTrace *string) {
	logCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_ = s.repo.Create(logCtx, errorlog.ErrorLogCreateInput{
		UserID:     userID,
		HTTPMethod: method,
		Route:      route,
		StatusCode: statusCode,
		Message:    errMsg,
		StackTrace: stackTrace,
	})
}
