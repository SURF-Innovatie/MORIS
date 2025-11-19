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
	title, description string,
	start, end time.Time,
	people []entities.Person,
	org entities.Organisation,
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
		Base:         base(id),
		Title:        title,
		Description:  description,
		StartDate:    start,
		EndDate:      end,
		People:       people,
		Organisation: org,
	}, nil
}

// ChangeTitle emits TitleChanged when different.
func ChangeTitle(id uuid.UUID, cur *entities.Project, title string) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if title == "" || cur.Title == title {
		return nil, nil
	}
	return events.TitleChanged{Base: base(id), Title: title}, nil
}

// ChangeDescription emits DescriptionChanged when different.
func ChangeDescription(id uuid.UUID, cur *entities.Project, desc string) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.Description == desc {
		return nil, nil
	}
	return events.DescriptionChanged{Base: base(id), Description: desc}, nil
}

// ChangeStartDate emits StartDateChanged when the start date differs and is valid.
func ChangeStartDate(id uuid.UUID, cur *entities.Project, start time.Time) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.StartDate.Equal(start) {
		return nil, nil
	}
	return events.StartDateChanged{Base: base(id), StartDate: start}, nil
}

// ChangeEndDate emits EndDateChanged when the end date differs and is valid.
func ChangeEndDate(id uuid.UUID, cur *entities.Project, end time.Time) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.EndDate.Equal(end) {
		return nil, nil
	}
	return events.EndDateChanged{Base: base(id), EndDate: end}, nil
}

// SetOrganisation emits OrganisationChanged when different.
func SetOrganisation(id uuid.UUID, cur *entities.Project, org entities.Organisation) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur.Organisation == org {
		return nil, nil
	}
	return events.OrganisationChanged{Base: base(id), Organisation: org}, nil
}

// AddPerson emits PersonAdded when not present.
func AddPerson(id uuid.UUID, cur *entities.Project, p entities.Person) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	for _, x := range cur.People {
		if x != nil && x.Name == p.Name {
			return nil, errors.New(fmt.Sprintf("person %s already exists in project %s", p.Id, cur.Id))
		}
	}
	return events.PersonAdded{Base: base(id), Person: p}, nil
}

// RemovePerson emits PersonRemoved when present.
func RemovePerson(id uuid.UUID, cur *entities.Project, p entities.Person) (events.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}

	exist := false
	for _, x := range cur.People {
		if x != nil && x.Name == p.Name {
			exist = true
		}
	}
	if !exist {
		return nil, errors.New(fmt.Sprintf("person %s not found for project %s", p.Id, cur.Id))
	}

	return events.PersonRemoved{Base: base(id), Person: p}, nil
}

func base(id uuid.UUID) events.Base {
	return events.Base{ProjectID: id, At: time.Now().UTC()}
}
