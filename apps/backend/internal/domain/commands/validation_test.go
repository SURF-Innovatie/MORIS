package commands_test

import (
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/google/uuid"
)

func Test_StartProject_Validation(t *testing.T) {
	id := uuid.New()
	start := time.Now().UTC()
	end := start.Add(-time.Hour)

	actor := uuid.New()
	_, err := commands.StartProject(id, actor, "", "d", start, start, nil, uuid.New())
	if err == nil {
		t.Fatal("missing title should error")
	}

	_, err = commands.StartProject(id, actor, "t", "d", start, end, nil, uuid.New())
	if err == nil {
		t.Fatal("end before start should error")
	}
}
