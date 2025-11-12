package commands_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

func Test_StartProject_Validation(t *testing.T) {
	id := uuid.New()
	start := time.Now().UTC()
	end := start.Add(-time.Hour)

	_, err := commands.StartProject(id, "", "d", start, start, nil, entities.Organisation{})
	if err == nil {
		t.Fatal("missing title should error")
	}

	_, err = commands.StartProject(id, "t", "d", start, end, nil, entities.Organisation{})
	if err == nil {
		t.Fatal("end before start should error")
	}
}
