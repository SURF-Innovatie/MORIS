package errorlog

import (
	"context"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// Service defines the interface for logging errors.
type Service interface {
	Log(ctx context.Context, userID, method, route string, statusCode int, errMsg, stackTrace string)
}

type service struct {
	client *ent.Client
}

// NewService creates a new error logging service.
func NewService(client *ent.Client) Service {
	return &service{client: client}
}

// Log records an error in the database.
func (s *service) Log(ctx context.Context, userID, method, route string, statusCode int, errMsg, stackTrace string) {
	// Use a detached context with timeout to ensure logging persists even if the request is canceled.
	logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	create := s.client.ErrorLog.Create().
		SetHTTPMethod(method).
		SetRoute(route).
		SetStatusCode(statusCode).
		SetErrorMessage(errMsg).
		SetNillableStackTrace(&stackTrace)

	if userID != "" {
		create.SetUserID(userID)
	}

	// Intentionally ignore error; logging failure shouldn't impact the main request flow.
	_ = create.Exec(logCtx)
}

// LogWithUUID is a helper if userID is passed as UUID
func (s *service) LogWithUUID(ctx context.Context, userID uuid.UUID, method, route string, statusCode int, errMsg, stackTrace string) {
	uidStr := ""
	if userID != uuid.Nil {
		uidStr = userID.String()
	}
	s.Log(ctx, uidStr, method, route, statusCode, errMsg, stackTrace)
}
