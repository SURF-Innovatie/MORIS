package raid

import (
	"context"
	"time"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// MapperService provides methods to map projects to RAiD requests and check compatibility.
type MapperService interface {
	// MapToCreateRequest builds a RAiD create request from a project's event stream.
	MapToCreateRequest(ctx context.Context, projectID uuid.UUID) (*raid.RAiDCreateRequest, error)

	// MapToUpdateRequest builds a RAiD update request (requires existing RAiD identifier).
	MapToUpdateRequest(ctx context.Context, projectID uuid.UUID, existingRAiD raid.RAiDId) (*raid.RAiDUpdateRequest, error)

	// CheckCompatibility validates a project against RAiD constraints.
	CheckCompatibility(ctx context.Context, projectID uuid.UUID) ([]Incompatibility, error)
}

// EventStore is the interface for loading project events.
type EventStore interface {
	Load(ctx context.Context, projectID uuid.UUID) ([]events.Event, int, error)
}

// PersonRepository is the interface for loading person entities.
type PersonRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Person, error)
}

// OrganisationRepository is the interface for loading organisation entities.
type OrganisationRepository interface {
	GetNode(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error)
}

// mapperService implements MapperService.
type mapperService struct {
	eventStore  EventStore
	persons     PersonRepository
	orgs        OrganisationRepository
	defaultLang string
}

// NewMapperService creates a new MapperService.
func NewMapperService(es EventStore, persons PersonRepository, orgs OrganisationRepository) MapperService {
	return &mapperService{
		eventStore:  es,
		persons:     persons,
		orgs:        orgs,
		defaultLang: raid.LanguageEnglish,
	}
}

func (s *mapperService) MapToCreateRequest(ctx context.Context, projectID uuid.UUID) (*raid.RAiDCreateRequest, error) {
	evts, _, err := s.eventStore.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project := s.reduceProject(evts)
	contributors, err := s.buildContributors(ctx, evts)
	if err != nil {
		return nil, err
	}
	orgs, err := s.buildOrganisations(ctx, project.OwningOrgNodeID, evts)
	if err != nil {
		return nil, err
	}

	return &raid.RAiDCreateRequest{
		Title:        s.buildTitles(evts),
		Date:         s.buildDate(project.StartDate, project.EndDate),
		Description:  s.buildDescriptions(evts),
		Access:       s.buildAccess(),
		Contributor:  contributors,
		Organisation: orgs,
		AlternateIdentifier: []raid.RAiDAlternateIdentifier{{
			Id:   projectID.String(),
			Type: "moris-project-id",
		}},
	}, nil
}

func (s *mapperService) MapToUpdateRequest(ctx context.Context, projectID uuid.UUID, existingRAiD raid.RAiDId) (*raid.RAiDUpdateRequest, error) {
	evts, _, err := s.eventStore.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project := s.reduceProject(evts)
	contributors, err := s.buildContributors(ctx, evts)
	if err != nil {
		return nil, err
	}
	orgs, err := s.buildOrganisations(ctx, project.OwningOrgNodeID, evts)
	if err != nil {
		return nil, err
	}

	return &raid.RAiDUpdateRequest{
		Identifier:   existingRAiD,
		Title:        s.buildTitles(evts),
		Date:         s.buildDate(project.StartDate, project.EndDate),
		Description:  s.buildDescriptions(evts),
		Access:       s.buildAccess(),
		Contributor:  contributors,
		Organisation: orgs,
		AlternateIdentifier: []raid.RAiDAlternateIdentifier{{
			Id:   projectID.String(),
			Type: "moris-project-id",
		}},
	}, nil
}

func (s *mapperService) CheckCompatibility(ctx context.Context, projectID uuid.UUID) ([]Incompatibility, error) {
	evts, _, err := s.eventStore.Load(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var incompatibilities []Incompatibility

	// Extract project state
	project := s.reduceProject(evts)
	titles := s.extractTitleHistory(evts)

	// Check: At least one current primary title
	activeTitles := lo.Filter(titles, func(t titleEntry, _ int) bool {
		now := time.Now()
		return t.startDate.Before(now) && (t.endDate == nil || t.endDate.After(now))
	})
	if len(activeTitles) == 0 {
		incompatibilities = append(incompatibilities, NewIncompatibility(NoActivePrimaryTitle))
	}
	if len(activeTitles) > 1 {
		incompatibilities = append(incompatibilities, NewIncompatibility(MultipleActivePrimaryTitle))
	}

	// Check: Title max length
	for _, t := range titles {
		if len(t.text) > raid.MaxTitleLength {
			incompatibilities = append(incompatibilities, NewIncompatibility(ProjectTitleTooLong))
		}
	}

	// Check: Description max length
	desc := s.extractDescription(evts)
	if len(desc) > raid.MaxDescriptionLength {
		incompatibilities = append(incompatibilities, NewIncompatibility(ProjectDescriptionTooLong))
	}

	// Check: Contributors
	members := s.extractMembers(evts)
	if len(members) == 0 {
		incompatibilities = append(incompatibilities, NewIncompatibility(NoContributors))
	}

	// Check contributors have ORCID
	for _, m := range members {
		person, err := s.persons.Get(ctx, m.personID)
		if err != nil || person == nil {
			continue
		}
		if lo.FromPtrOr(person.ORCiD, "") == "" {
			incompatibilities = append(incompatibilities, NewIncompatibility(ContributorWithoutOrcid).WithObjectID(m.personID))
		}
	}

	// Check: Organisation has ROR
	if project.OwningOrgNodeID != uuid.Nil {
		org, err := s.orgs.GetNode(ctx, project.OwningOrgNodeID)
		if err == nil && org != nil {
			if lo.FromPtrOr(org.RorID, "") == "" {
				incompatibilities = append(incompatibilities, NewIncompatibility(OrganisationWithoutRor).WithObjectID(org.ID))
			}
		}
	} else {
		incompatibilities = append(incompatibilities, NewIncompatibility(NoLeadResearchOrganisation))
	}

	return incompatibilities, nil
}

// --- Helper types ---

type titleEntry struct {
	text      string
	startDate time.Time
	endDate   *time.Time
}

type memberEntry struct {
	personID  uuid.UUID
	roleID    uuid.UUID
	startDate time.Time
	endDate   *time.Time
}

type projectState struct {
	StartDate       time.Time
	EndDate         time.Time
	Title           string
	Description     string
	OwningOrgNodeID uuid.UUID
}

// --- Reduction helpers ---

func (s *mapperService) reduceProject(evts []events.Event) projectState {
	var p projectState
	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			p.StartDate = evt.StartDate
			p.EndDate = evt.EndDate
			p.Title = evt.Title
			p.Description = evt.Description
			p.OwningOrgNodeID = evt.OwningOrgNodeID
		case *events.TitleChanged:
			p.Title = evt.Title
		case *events.DescriptionChanged:
			p.Description = evt.Description
		case *events.StartDateChanged:
			p.StartDate = evt.StartDate
		case *events.EndDateChanged:
			p.EndDate = evt.EndDate
		case *events.OwningOrgNodeChanged:
			p.OwningOrgNodeID = evt.OwningOrgNodeID
		}
	}
	return p
}

func (s *mapperService) extractTitleHistory(evts []events.Event) []titleEntry {
	var titles []titleEntry
	var current *titleEntry

	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			if current != nil {
				current.endDate = lo.ToPtr(evt.OccurredAt())
			}
			titles = append(titles, titleEntry{text: evt.Title, startDate: evt.OccurredAt()})
			current = &titles[len(titles)-1]
		case *events.TitleChanged:
			if current != nil {
				current.endDate = lo.ToPtr(evt.OccurredAt())
			}
			titles = append(titles, titleEntry{text: evt.Title, startDate: evt.OccurredAt()})
			current = &titles[len(titles)-1]
		}
	}
	return titles
}

func (s *mapperService) extractDescription(evts []events.Event) string {
	desc := ""
	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			desc = evt.Description
		case *events.DescriptionChanged:
			desc = evt.Description
		}
	}
	return desc
}

func (s *mapperService) extractMembers(evts []events.Event) []memberEntry {
	members := make(map[uuid.UUID]*memberEntry)
	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectRoleAssigned:
			members[evt.PersonID] = &memberEntry{
				personID:  evt.PersonID,
				roleID:    evt.ProjectRoleID,
				startDate: evt.OccurredAt(),
			}
		case *events.ProjectRoleUnassigned:
			if m, ok := members[evt.PersonID]; ok {
				m.endDate = lo.ToPtr(evt.OccurredAt())
			}
		}
	}
	out := make([]memberEntry, 0, len(members))
	for _, m := range members {
		out = append(out, *m)
	}
	return out
}

// --- Build helpers ---

func (s *mapperService) buildTitles(evts []events.Event) []raid.RAiDTitle {
	history := s.extractTitleHistory(evts)
	return lo.Map(history, func(t titleEntry, _ int) raid.RAiDTitle {
		return raid.RAiDTitle{
			Text:      t.text,
			Type:      raid.RAiDTitleType{Id: raid.TitleTypePrimaryURI, SchemaUri: raid.TitleTypeSchemaURI},
			StartDate: t.startDate.Format("2006-01-02"),
			EndDate:   lo.TernaryF(t.endDate != nil, func() *string { return lo.ToPtr(t.endDate.Format("2006-01-02")) }, func() *string { return nil }),
			Language:  &raid.RAiDLanguage{Id: s.defaultLang, SchemaUri: raid.LanguageSchemaURI},
		}
	})
}

func (s *mapperService) buildDescriptions(evts []events.Event) []raid.RAiDDescription {
	desc := s.extractDescription(evts)
	if desc == "" {
		return nil
	}
	return []raid.RAiDDescription{{
		Text:     desc,
		Type:     raid.RAiDDescriptionType{Id: raid.DescriptionTypePrimaryURI, SchemaUri: raid.DescriptionTypeSchemaURI},
		Language: &raid.RAiDLanguage{Id: s.defaultLang, SchemaUri: raid.LanguageSchemaURI},
	}}
}

func (s *mapperService) buildDate(start, end time.Time) *raid.RAiDDate {
	if start.IsZero() {
		return nil
	}
	d := &raid.RAiDDate{StartDate: start.Format("2006-01-02")}
	if !end.IsZero() {
		d.EndDate = lo.ToPtr(end.Format("2006-01-02"))
	}
	return d
}

func (s *mapperService) buildAccess() raid.RAiDAccess {
	return raid.RAiDAccess{
		Type: raid.RAiDAccessType{Id: raid.AccessTypeOpenURI, SchemaUri: raid.AccessTypeSchemaURI},
	}
}

func (s *mapperService) buildContributors(ctx context.Context, evts []events.Event) ([]raid.RAiDContributor, error) {
	members := s.extractMembers(evts)
	personIDs := lo.Map(members, func(m memberEntry, _ int) uuid.UUID { return m.personID })
	persons, err := s.persons.GetByIDs(ctx, personIDs)
	if err != nil {
		return nil, err
	}
	personMap := lo.KeyBy(persons, func(p entities.Person) uuid.UUID { return p.ID })

	return lo.FilterMap(members, func(m memberEntry, _ int) (raid.RAiDContributor, bool) {
		p, ok := personMap[m.personID]
		if !ok || lo.FromPtrOr(p.ORCiD, "") == "" {
			return raid.RAiDContributor{}, false
		}
		return raid.RAiDContributor{
			Id:        p.ORCiD,
			SchemaUri: raid.ContributorSchemaORCID,
			Position: []raid.RAiDContributorPosition{{
				Id:        raid.ContributorPositionOtherURI,
				SchemaUri: raid.ContributorPositionSchemaURI,
				StartDate: m.startDate.Format("2006-01-02"),
				EndDate:   lo.TernaryF(m.endDate != nil, func() *string { return lo.ToPtr(m.endDate.Format("2006-01-02")) }, func() *string { return nil }),
			}},
			Role: []raid.RAiDContributorRole{{
				Id:        raid.ContributorRoleProjectAdminURI,
				SchemaUri: raid.ContributorRoleSchemaURI,
			}},
		}, true
	}), nil
}

func (s *mapperService) buildOrganisations(ctx context.Context, orgID uuid.UUID, evts []events.Event) ([]raid.RAiDOrganisation, error) {
	if orgID == uuid.Nil {
		return nil, nil
	}
	org, err := s.orgs.GetNode(ctx, orgID)
	if err != nil || org == nil || lo.FromPtrOr(org.RorID, "") == "" {
		return nil, nil
	}

	// Find when the org was set
	startDate := time.Now()
	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			startDate = evt.OccurredAt()
		case *events.OwningOrgNodeChanged:
			if evt.OwningOrgNodeID == orgID {
				startDate = evt.OccurredAt()
			}
		}
	}

	return []raid.RAiDOrganisation{{
		Id:        *org.RorID,
		SchemaUri: raid.OrganisationSchemaROR,
		Role: []raid.RAiDOrganisationRole{{
			Id:        raid.OrganisationRoleLeadResearchURI,
			SchemaUri: raid.OrganisationRoleSchemaURI,
			StartDate: startDate.Format("2006-01-02"),
		}},
	}}, nil
}
