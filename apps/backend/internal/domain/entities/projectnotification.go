package entities

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

type ProjectNotification struct {
	Id        uuid.UUID
	User      *ent.User
	Message   string
	SentAt    time.Time
	ProjectId uuid.UUID
}
