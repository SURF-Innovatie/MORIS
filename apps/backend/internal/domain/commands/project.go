package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// StartProject creates the initial snapshot event.
func StartProject(
	id uuid.UUID,
	actor uuid.UUID,
	title, description string,
	start, end time.Time,
	members []entities.ProjectMember,
	org uuid.UUID,
) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if end.Before(start) {
		return nil, errors.New("end date before start date")
	}
	return &events.ProjectStarted{
		Base:            base(id, actor, event.StatusApproved),
		Title:           title,
		Description:     description,
		StartDate:       start,
		EndDate:         end,
		Members:         members,
		OwningOrgNodeID: org,
	}, nil
}

// ChangeTitle emits TitleChanged when different.
func ChangeTitle(id uuid.UUID, actor uuid.UUID, cur *entities.Project, title string, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if title == "" || cur.Title == title {
		return nil, nil
	}
	return &events.TitleChanged{Base: base(id, actor, status), Title: title}, nil
}

// ChangeDescription emits DescriptionChanged when different.
func ChangeDescription(id uuid.UUID, actor uuid.UUID, cur *entities.Project, desc string, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.Description == desc {
		return nil, nil
	}
	return &events.DescriptionChanged{Base: base(id, actor, status), Description: desc}, nil
}

// ChangeStartDate emits StartDateChanged when the start date differs and is valid.
func ChangeStartDate(id uuid.UUID, actor uuid.UUID, cur *entities.Project, start time.Time, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.StartDate.Equal(start) {
		return nil, nil
	}
	return &events.StartDateChanged{Base: base(id, actor, status), StartDate: start}, nil
}

// ChangeEndDate emits EndDateChanged when the end date differs and is valid.
func ChangeEndDate(id uuid.UUID, actor uuid.UUID, cur *entities.Project, end time.Time, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.EndDate.Equal(end) {
		return nil, nil
	}
	return &events.EndDateChanged{Base: base(id, actor, status), EndDate: end}, nil
}

func ChangeOwningOrgNode(
	id uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	orgNodeID uuid.UUID,
	status event.Status,
) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if orgNodeID == uuid.Nil {
		return nil, errors.New("organisation node id is required")
	}
	if cur.OwningOrgNodeID == orgNodeID {
		return nil, nil
	}

	return &events.OwningOrgNodeChanged{
		Base:            base(id, actor, status),
		OwningOrgNodeID: orgNodeID,
	}, nil
}

func AssignProjectRole(
	id uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	personID uuid.UUID,
	projectRoleID uuid.UUID,
	status event.Status,
) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if personID == uuid.Nil {
		return nil, errors.New("person id is required")
	}

	// idempotency: don't assign same role twice
	for _, m := range cur.Members {
		if m.PersonID == personID && m.ProjectRoleID == projectRoleID {
			return nil, nil
		}
	}

	return &events.ProjectRoleAssigned{
		Base:          base(id, actor, status),
		PersonID:      personID,
		ProjectRoleID: projectRoleID,
	}, nil
}

func UnassignProjectRole(
	id uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	personID uuid.UUID,
	projectRoleID uuid.UUID,
	status event.Status,
) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if personID == uuid.Nil {
		return nil, errors.New("person id is required")
	}

	found := false
	for _, m := range cur.Members {
		if m.PersonID == personID && m.ProjectRoleID == projectRoleID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil // idempotent remove
	}

	return &events.ProjectRoleUnassigned{
		Base:          base(id, actor, status),
		PersonID:      personID,
		ProjectRoleID: projectRoleID,
	}, nil
}

func SetOwningOrgNode(
	id uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	owningOrgNodeID uuid.UUID,
	status event.Status,
) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if owningOrgNodeID == uuid.Nil {
		return nil, errors.New("owning org node id is required")
	}
	if cur.OwningOrgNodeID == owningOrgNodeID {
		return nil, nil
	}
	return &events.OwningOrgNodeChanged{
		Base:            base(id, actor, status),
		OwningOrgNodeID: owningOrgNodeID,
	}, nil
}

// AddProduct emits ProductAdded when not present
func AddProduct(id uuid.UUID, actor uuid.UUID, cur *entities.Project, productID uuid.UUID, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}

	for _, x := range cur.ProductIDs {
		if x == productID {
			return nil, fmt.Errorf("product %s already exists in project %s", productID, cur.Id)
		}
	}

	return &events.ProductAdded{Base: base(id, actor, status), ProductID: productID}, nil
}

// RemoveProduct emits ProductRemoved when present
func RemoveProduct(id uuid.UUID, actor uuid.UUID, cur *entities.Project, productID uuid.UUID, status event.Status) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	exist := false
	for _, x := range cur.ProductIDs {
		if x == productID {
			exist = true
		}
	}
	if !exist {
		return nil, fmt.Errorf("product %s not found for project %s", productID, cur.Id)
	}
	return &events.ProductRemoved{Base: base(id, actor, status), ProductID: productID}, nil
}

func base(id uuid.UUID, actor uuid.UUID, status event.Status) events.Base {
	return events.Base{ProjectID: id, At: time.Now().UTC(), ID: uuid.New(), CreatedBy: actor, Status: string(status)}
}
