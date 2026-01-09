package raidsink

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

// Schema URIs for RAiD metadata
const (
	// Title types - https://vocabulary.raid.org/title.type.schema/376
	TitleTypePrimaryURI = "https://vocabulary.raid.org/title.type.schema/5"
	TitleTypeSchemaURI  = "https://vocabulary.raid.org/title.type.schema/376"

	// Description types - https://vocabulary.raid.org/description.type.schema/329
	DescriptionTypePrimaryURI = "https://vocabulary.raid.org/description.type.schema/318"
	DescriptionTypeSchemaURI  = "https://vocabulary.raid.org/description.type.schema/329"

	// Access types - https://vocabulary.raid.org/access.type.schema/289
	AccessTypeOpenURI   = "https://vocabulary.raid.org/access.type.schema/238"
	AccessTypeSchemaURI = "https://vocabulary.raid.org/access.type.schema/289"

	// Language - ISO 639-3
	LanguageSchemaURI = "https://www.iso.org/standard/39534.html"
	LanguageEnglish   = "eng"

	// Contributor schemas
	ContributorSchemaORCID = "https://orcid.org/"

	// Contributor position - https://vocabulary.raid.org/contributor.position.schema/305
	ContributorPositionOtherURI  = "https://vocabulary.raid.org/contributor.position.schema/307"
	ContributorPositionSchemaURI = "https://vocabulary.raid.org/contributor.position.schema/305"

	// Contributor role (CRediT) - https://credit.niso.org/
	ContributorRoleSchemaURI       = "https://credit.niso.org/"
	ContributorRoleProjectAdminURI = "https://credit.niso.org/contributor-roles/project-administration/"

	// Organisation schemas
	OrganisationSchemaROR = "https://ror.org/"

	// Organisation role - https://vocabulary.raid.org/organisation.role.schema/359
	OrganisationRoleLeadURI   = "https://vocabulary.raid.org/organisation.role.schema/182"
	OrganisationRoleSchemaURI = "https://vocabulary.raid.org/organisation.role.schema/359"
)

// titleEntry represents a historical title with date range
type titleEntry struct {
	text      string
	startDate time.Time
	endDate   *time.Time
}

// contributorState tracks contributor state from event processing
type contributorState struct {
	personID  uuid.UUID
	person    *entities.Person
	roleID    uuid.UUID
	startDate time.Time
	endDate   *time.Time
}

type memberKey struct {
	personID uuid.UUID
	roleID   uuid.UUID
}

// RAiDMapper handles ProjectContext -> RAiD transformation by processing events.
type RAiDMapper struct {
	defaultLanguage string
}

// RAiDMapperOption configures the RAiD mapper.
type RAiDMapperOption func(*RAiDMapper)

// WithDefaultLanguage sets the default language for RAiD metadata.
func WithDefaultLanguage(lang string) RAiDMapperOption {
	return func(m *RAiDMapper) {
		m.defaultLanguage = lang
	}
}

// NewRAiDMapper creates a new RAiD mapper with the given options.
func NewRAiDMapper(opts ...RAiDMapperOption) *RAiDMapper {
	m := &RAiDMapper{defaultLanguage: LanguageEnglish}
	lo.ForEach(opts, func(opt RAiDMapperOption, _ int) { opt(m) })
	return m
}

// MapToCreateRequest converts a ProjectContext to a RAiD create request.
func (m *RAiDMapper) MapToCreateRequest(pc adapter.ProjectContext) *raid.RAiDCreateRequest {
	titles := m.extractTitleHistory(pc.Events)
	descriptions := m.extractDescriptionHistory(pc.Events)
	startDate, endDate := m.extractDates(pc.Events)

	return &raid.RAiDCreateRequest{
		Title:        m.mapTitles(titles),
		Description:  m.mapDescriptions(descriptions),
		Date:         m.mapDate(startDate, endDate),
		Access:       m.defaultAccess(),
		Contributor:  m.extractContributors(pc.Events, pc.Members),
		Organisation: m.extractOrganisations(pc.Events, pc.OrgNode),
	}
}

// extractTitleHistory builds a list of all titles with their date ranges
func (m *RAiDMapper) extractTitleHistory(evts []events.Event) []titleEntry {
	var titles []titleEntry
	var currentTitle *titleEntry

	endCurrentTitle := func(at time.Time) {
		if currentTitle != nil {
			currentTitle.endDate = lo.ToPtr(at)
		}
	}

	startNewTitle := func(text string, at time.Time) {
		titles = append(titles, titleEntry{text: text, startDate: at})
		currentTitle = &titles[len(titles)-1]
	}

	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			endCurrentTitle(evt.OccurredAt())
			startNewTitle(evt.Title, evt.OccurredAt())
		case *events.TitleChanged:
			endCurrentTitle(evt.OccurredAt())
			startNewTitle(evt.Title, evt.OccurredAt())
		}
	}

	return titles
}

// extractDescriptionHistory returns the current description
func (m *RAiDMapper) extractDescriptionHistory(evts []events.Event) []string {
	description := lo.Reduce(evts, func(acc string, e events.Event, _ int) string {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			return evt.Description
		case *events.DescriptionChanged:
			return evt.Description
		}
		return acc
	}, "")

	return lo.Ternary(description == "", nil, []string{description})
}

// extractDates finds start/end dates from the event stream
func (m *RAiDMapper) extractDates(evts []events.Event) (time.Time, *time.Time) {
	var startDate, endDate time.Time

	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			startDate, endDate = evt.StartDate, evt.EndDate
		case *events.StartDateChanged:
			startDate = evt.StartDate
		case *events.EndDateChanged:
			endDate = evt.EndDate
		}
	}

	return startDate, lo.Ternary(endDate.IsZero(), nil, &endDate)
}

// extractContributors processes role assignment events to build contributor list
func (m *RAiDMapper) extractContributors(evts []events.Event, persons []entities.Person) []raid.RAiDContributor {
	personMap := lo.KeyBy(persons, func(p entities.Person) uuid.UUID { return p.ID })
	states := make(map[memberKey]*contributorState)

	for _, e := range evts {
		switch evt := e.(type) {
		case *events.ProjectRoleAssigned:
			key := memberKey{personID: evt.PersonID, roleID: evt.ProjectRoleID}
			person := lo.ToPtr(personMap[evt.PersonID])
			states[key] = &contributorState{
				personID:  evt.PersonID,
				person:    person,
				roleID:    evt.ProjectRoleID,
				startDate: evt.OccurredAt(),
			}
		case *events.ProjectRoleUnassigned:
			key := memberKey{personID: evt.PersonID, roleID: evt.ProjectRoleID}
			if state, ok := states[key]; ok {
				state.endDate = lo.ToPtr(evt.OccurredAt())
			}
		}
	}

	// Filter to only those with ORCID and map to RAiD contributors
	return lo.FilterMap(lo.Values(states), func(s *contributorState, _ int) (raid.RAiDContributor, bool) {
		if s.person == nil || lo.FromPtrOr(s.person.ORCiD, "") == "" {
			return raid.RAiDContributor{}, false
		}

		position := raid.RAiDContributorPosition{
			Id:        ContributorPositionOtherURI,
			SchemaUri: ContributorPositionSchemaURI,
			StartDate: s.startDate.Format("2006-01-02"),
			EndDate:   lo.TernaryF(s.endDate != nil, func() *string { return lo.ToPtr(s.endDate.Format("2006-01-02")) }, func() *string { return nil }),
		}

		return raid.RAiDContributor{
			Id:        s.person.ORCiD,
			SchemaUri: ContributorSchemaORCID,
			Position:  []raid.RAiDContributorPosition{position},
			Role: []raid.RAiDContributorRole{{
				Id:        ContributorRoleProjectAdminURI,
				SchemaUri: ContributorRoleSchemaURI,
			}},
		}, true
	})
}

// extractOrganisations finds organisations with roles from events
func (m *RAiDMapper) extractOrganisations(evts []events.Event, orgNode *entities.OrganisationNode) []raid.RAiDOrganisation {
	if orgNode == nil || lo.FromPtrOr(orgNode.RorID, "") == "" {
		return nil
	}

	startDate := lo.Reduce(evts, func(acc time.Time, e events.Event, _ int) time.Time {
		switch evt := e.(type) {
		case *events.ProjectStarted:
			return evt.OccurredAt()
		case *events.OwningOrgNodeChanged:
			return evt.OccurredAt()
		}
		return acc
	}, time.Now())

	return []raid.RAiDOrganisation{{
		Id:        *orgNode.RorID,
		SchemaUri: OrganisationSchemaROR,
		Role: []raid.RAiDOrganisationRole{{
			Id:        OrganisationRoleLeadURI,
			SchemaUri: OrganisationRoleSchemaURI,
			StartDate: startDate.Format("2006-01-02"),
		}},
	}}
}

func (m *RAiDMapper) mapTitles(titles []titleEntry) []raid.RAiDTitle {
	return lo.Map(titles, func(t titleEntry, _ int) raid.RAiDTitle {
		return raid.RAiDTitle{
			Text:      t.text,
			Type:      raid.RAiDTitleType{Id: TitleTypePrimaryURI, SchemaUri: TitleTypeSchemaURI},
			StartDate: t.startDate.Format("2006-01-02"),
			EndDate:   lo.TernaryF(t.endDate != nil, func() *string { return lo.ToPtr(t.endDate.Format("2006-01-02")) }, func() *string { return nil }),
			Language:  &raid.RAiDLanguage{Id: m.defaultLanguage, SchemaUri: LanguageSchemaURI},
		}
	})
}

func (m *RAiDMapper) mapDescriptions(descriptions []string) []raid.RAiDDescription {
	return lo.FilterMap(descriptions, func(d string, _ int) (raid.RAiDDescription, bool) {
		if d == "" {
			return raid.RAiDDescription{}, false
		}
		return raid.RAiDDescription{
			Text:     d,
			Type:     raid.RAiDDescriptionType{Id: DescriptionTypePrimaryURI, SchemaUri: DescriptionTypeSchemaURI},
			Language: &raid.RAiDLanguage{Id: m.defaultLanguage, SchemaUri: LanguageSchemaURI},
		}, true
	})
}

func (m *RAiDMapper) mapDate(startDate time.Time, endDate *time.Time) *raid.RAiDDate {
	if startDate.IsZero() {
		return nil
	}
	return &raid.RAiDDate{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   lo.TernaryF(endDate != nil && !endDate.IsZero(), func() *string { return lo.ToPtr(endDate.Format("2006-01-02")) }, func() *string { return nil }),
	}
}

func (m *RAiDMapper) defaultAccess() raid.RAiDAccess {
	return raid.RAiDAccess{
		Type: raid.RAiDAccessType{Id: AccessTypeOpenURI, SchemaUri: AccessTypeSchemaURI},
	}
}
