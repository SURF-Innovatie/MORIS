package errorlog

import (
	"time"

	"github.com/google/uuid"
)

type ErrorLog struct {
	ID         uuid.UUID
	UserID     *uuid.UUID
	HTTPMethod string
	Route      string
	StatusCode int
	Message    string
	StackTrace *string
	CreatedAt  time.Time
}

type ErrorLogCreateInput struct {
	UserID     *uuid.UUID
	HTTPMethod string
	Route      string
	StatusCode int
	Message    string
	StackTrace *string
}
