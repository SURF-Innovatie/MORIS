package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const ActualRecordedType = "budget.actual_recorded"

// ActualRecorded represents the event when an actual expenditure is recorded
type ActualRecorded struct {
	Base
	ActualID     uuid.UUID `json:"actualId"`
	LineItemID   uuid.UUID `json:"lineItemId"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	RecordedDate time.Time `json:"recordedDate"`
	Source       string    `json:"source"` // "manual" | "erp_sync"
}

func (ActualRecorded) isEvent()     {}
func (ActualRecorded) Type() string { return ActualRecordedType }
func (e ActualRecorded) String() string {
	return fmt.Sprintf("Actual recorded: €%.2f on %s", e.Amount, e.RecordedDate.Format("2006-01-02"))
}

func (e *ActualRecorded) NotificationMessage() string {
	return fmt.Sprintf("Expenditure of €%.2f has been recorded.", e.Amount)
}

type ActualRecordedInput struct {
	LineItemID   uuid.UUID `json:"lineItemId"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	RecordedDate time.Time `json:"recordedDate"`
	Source       string    `json:"source"` // "manual" | "erp_sync"
}

func DecideActualRecorded(
	projectID uuid.UUID,
	actor uuid.UUID,
	in ActualRecordedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.LineItemID == uuid.Nil {
		return nil, errors.New("line item id is required")
	}
	if in.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	source := in.Source
	if source == "" {
		source = "manual"
	}

	recordedDate := in.RecordedDate
	if recordedDate.IsZero() {
		recordedDate = time.Now().UTC()
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ActualRecordedMeta.FriendlyName

	return &ActualRecorded{
		Base:         base,
		ActualID:     uuid.New(),
		LineItemID:   in.LineItemID,
		Amount:       in.Amount,
		Description:  in.Description,
		RecordedDate: recordedDate,
		Source:       source,
	}, nil
}

var ActualRecordedMeta = EventMeta{
	Type:         ActualRecordedType,
	FriendlyName: "Actual Recorded",
}

func init() {
	RegisterMeta(ActualRecordedMeta, func() Event {
		return &ActualRecorded{
			Base: Base{FriendlyNameStr: ActualRecordedMeta.FriendlyName},
		}
	})

	RegisterDecider[ActualRecordedInput](ActualRecordedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in ActualRecordedInput, status Status) (Event, error) {
			return DecideActualRecorded(projectID, actor, in, status)
		})

	RegisterInputType(ActualRecordedType, ActualRecordedInput{})
}
