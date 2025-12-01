package commands

import (
	"errors"
	"fmt"
	"time"

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
	people []uuid.UUID,
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
	return events.ProjectStarted{
		Base:           base(id, actor),
		Title:          title,
		Description:    description,
		StartDate:      start,
		EndDate:        end,
		People:         people,
		OrganisationID: org,
	}, nil
}

// ChangeTitle emits TitleChanged when different.
func ChangeTitle(id uuid.UUID, actor uuid.UUID, cur *entities.Project, title string) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if title == "" || cur.Title == title {
		return nil, nil
	}
	return events.TitleChanged{Base: base(id, actor), Title: title}, nil
}

// ChangeDescription emits DescriptionChanged when different.
func ChangeDescription(id uuid.UUID, actor uuid.UUID, cur *entities.Project, desc string) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.Description == desc {
		return nil, nil
	}
	return events.DescriptionChanged{Base: base(id, actor), Description: desc}, nil
}

// ChangeStartDate emits StartDateChanged when the start date differs and is valid.
func ChangeStartDate(id uuid.UUID, actor uuid.UUID, cur *entities.Project, start time.Time) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.StartDate.Equal(start) {
		return nil, nil
	}
	return events.StartDateChanged{Base: base(id, actor), StartDate: start}, nil
}

// ChangeEndDate emits EndDateChanged when the end date differs and is valid.
func ChangeEndDate(id uuid.UUID, actor uuid.UUID, cur *entities.Project, end time.Time) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.EndDate.Equal(end) {
		return nil, nil
	}
	return events.EndDateChanged{Base: base(id, actor), EndDate: end}, nil
}

// SetOrganisation emits OrganisationChanged when different.
func SetOrganisation(id uuid.UUID, actor uuid.UUID, cur *entities.Project, org uuid.UUID) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.Organisation == org {
		return nil, nil
	}
	return events.OrganisationChanged{Base: base(id, actor), OrganisationID: org}, nil
}

// AddPerson emits PersonAdded when not present.
func AddPerson(id uuid.UUID, actor uuid.UUID, cur *entities.Project, personId uuid.UUID) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	for _, x := range cur.People {
		if x == personId {
			return nil, fmt.Errorf("person %s already exists in project %s", personId, cur.Id)
		}
	}
	return events.PersonAdded{Base: base(id, actor), PersonId: personId}, nil
}

// RemovePerson emits PersonRemoved when present.
func RemovePerson(id uuid.UUID, actor uuid.UUID, cur *entities.Project, personId uuid.UUID) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}

	exist := false
	for _, x := range cur.People {
		if x == personId {
			exist = true
		}
	}
	if !exist {
		return nil, fmt.Errorf("person %s not found for project %s", personId, cur.Id)
	}

	return events.PersonRemoved{Base: base(id, actor), PersonId: personId}, nil
}

// AddProduct emits ProductAdded when not present
func AddProduct(id uuid.UUID, actor uuid.UUID, cur *entities.Project, productID uuid.UUID) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}

	for _, x := range cur.Products {
		if x == productID {
			return nil, fmt.Errorf("product %s already exists in project %s", productID, cur.Id)
		}
	}
	return events.ProductAdded{Base: base(id, actor), ProductID: productID}, nil
}

// RemoveProduct emits ProductRemoved when present
func RemoveProduct(id uuid.UUID, actor uuid.UUID, cur *entities.Project, productID uuid.UUID) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	exist := false
	for _, x := range cur.People {
		if x == productID {
			exist = true
		}
	}
	if !exist {
		return nil, fmt.Errorf("product %s not found for project %s", productID, cur.Id)
	}
	return events.ProductRemoved{Base: base(id, actor), ProductID: productID}, nil
}

func base(id uuid.UUID, actor uuid.UUID) events.Base {
	return events.Base{ProjectID: id, At: time.Now().UTC(), ID: uuid.New(), CreatedBy: actor}
}
