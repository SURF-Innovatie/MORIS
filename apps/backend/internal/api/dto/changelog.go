package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type ChangelogEntry struct {
	Event string    `json:"event"`
	At    time.Time `json:"at"`
}

type Changelog struct {
	Entries []ChangelogEntry `json:"entries"`
}

func (c Changelog) FromEntity(log entities.ChangeLog) Changelog {
	entries := make([]ChangelogEntry, 0, len(log.Entries))
	for _, e := range log.Entries {
		entries = append(entries, ChangelogEntry{
			Event: e.Event,
			At:    e.At,
		})
	}
	return Changelog{Entries: entries}
}
