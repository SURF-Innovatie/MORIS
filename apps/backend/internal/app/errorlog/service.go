package errorlog

import (
	"context"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Log(ctx context.Context, userID *uuid.UUID, method, route string, statusCode int, errMsg string, stackTrace *string)
}

type service struct {
	repo    Repository
	timeout time.Duration
}

type Option func(*service)

func WithTimeout(d time.Duration) Option {
	return func(s *service) { s.timeout = d }
}

func NewService(repo Repository, opts ...Option) Service {
	s := &service{
		repo:    repo,
		timeout: 5 * time.Second,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *service) Log(_ context.Context, userID *uuid.UUID, method, route string, statusCode int, errMsg string, stackTrace *string) {
	// Detached context so logging still works if request ctx is cancelled.
	logCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_ = s.repo.Create(logCtx, entities.ErrorLogCreateInput{
		UserID:     userID,
		HTTPMethod: method,
		Route:      route,
		StatusCode: statusCode,
		Message:    errMsg,
		StackTrace: stackTrace,
	})
}
