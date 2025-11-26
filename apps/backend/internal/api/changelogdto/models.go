package changelogdto

import "time"

type ChangelogEntry struct {
	Event string    `json:"event"`
	At    time.Time `json:"at"`
}

type Changelog struct {
	Entries []ChangelogEntry `json:"entries"`
}
