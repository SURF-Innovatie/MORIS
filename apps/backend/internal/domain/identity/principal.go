package identity

import "github.com/google/uuid"

// Principal represents the authenticated actor for the current request.
type Principal struct {
	UserID     uuid.UUID
	PersonID   uuid.UUID
	IsSysAdmin bool
}
