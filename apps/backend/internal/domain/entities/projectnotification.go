package entities

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type Notification struct {
	Id      uuid.UUID
	User    *ent.User
	Event   *ent.Event
	Message string
	Read    bool
	SentAt  time.Time
}
