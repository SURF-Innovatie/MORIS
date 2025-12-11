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
	Type    string
	Read    bool
	SentAt  time.Time
}

func (n *Notification) FromEnt(row *ent.Notification, u *ent.User, e *ent.Event) *Notification {
	return &Notification{
		Id:      row.ID,
		User:    u,
		Event:   e,
		Message: row.Message,
		Type:    row.Type.String(),
		Read:    row.Read,
		SentAt:  row.SentAt,
	}
}
