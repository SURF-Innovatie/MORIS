package orcid_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/SURF-Innovatie/MORIS/ent"
	ext "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type fakeOrcidClient struct {
	authURL     string
	authErr     error
	exchangeOut string
	exchangeErr error
	searchOut   []ext.OrcidPerson
	searchErr   error
}

func (f *fakeOrcidClient) AuthURL() (string, error) {
	return f.authURL, f.authErr
}

func (f *fakeOrcidClient) ExchangeCode(context.Context, string) (string, error) {
	if f.exchangeErr != nil {
		return "", f.exchangeErr
	}
	return f.exchangeOut, nil
}

func (f *fakeOrcidClient) SearchExpanded(context.Context, string) ([]ext.OrcidPerson, error) {
	if f.searchErr != nil {
		return nil, f.searchErr
	}
	return f.searchOut, nil
}

// Fake user repo (implements internal/app/user.Repository)
type fakeUserRepo struct {
	byID       map[uuid.UUID]*ent.User
	byPersonID map[uuid.UUID]uuid.UUID
}

func (r *fakeUserRepo) Get(_ context.Context, id uuid.UUID) (*entities.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return transform.ToEntityPtr[entities.User](u), nil
}

func (r *fakeUserRepo) Create(_ context.Context, u entities.User) (*entities.User, error) {
	if r.byID == nil {
		r.byID = map[uuid.UUID]*ent.User{}
	}
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	eu := &ent.User{
		ID:         u.ID,
		PersonID:   u.PersonID,
		Password:   u.Password,
		IsSysAdmin: u.IsSysAdmin,
		IsActive:   u.IsActive,
	}
	if u.ZenodoAccessToken != nil {
		eu.ZenodoAccessToken = *u.ZenodoAccessToken
	}
	if u.ZenodoRefreshToken != nil {
		eu.ZenodoRefreshToken = *u.ZenodoRefreshToken
	}
	r.byID[u.ID] = eu
	if r.byPersonID == nil {
		r.byPersonID = map[uuid.UUID]uuid.UUID{}
	}
	r.byPersonID[u.PersonID] = u.ID
	return transform.ToEntityPtr[entities.User](eu), nil
}

func (r *fakeUserRepo) Update(_ context.Context, id uuid.UUID, u entities.User) (*entities.User, error) {
	if r.byID == nil {
		r.byID = map[uuid.UUID]*ent.User{}
	}
	eu := &ent.User{
		ID:         id,
		PersonID:   u.PersonID,
		Password:   u.Password,
		IsSysAdmin: u.IsSysAdmin,
		IsActive:   u.IsActive,
	}
	if u.ZenodoAccessToken != nil {
		eu.ZenodoAccessToken = *u.ZenodoAccessToken
	}
	if u.ZenodoRefreshToken != nil {
		eu.ZenodoRefreshToken = *u.ZenodoRefreshToken
	}
	r.byID[id] = eu
	if r.byPersonID == nil {
		r.byPersonID = map[uuid.UUID]uuid.UUID{}
	}
	r.byPersonID[u.PersonID] = id
	return transform.ToEntityPtr[entities.User](eu), nil
}

func (r *fakeUserRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.byID, id)
	return nil
}

func (r *fakeUserRepo) ToggleActive(_ context.Context, id uuid.UUID, isActive bool) error {
	u, ok := r.byID[id]
	if !ok {
		return errors.New("not found")
	}
	u.IsActive = isActive
	r.byID[id] = u
	return nil
}

func (r *fakeUserRepo) ListUsers(context.Context, int, int) ([]entities.User, int, error) {
	out := make([]*ent.User, 0, len(r.byID))
	for _, u := range r.byID {
		out = append(out, u)
	}
	return transform.ToEntities[entities.User](out), len(out), nil
}

func (r *fakeUserRepo) GetByPersonID(ctx context.Context, personID uuid.UUID) (*entities.User, error) {
	if r.byPersonID == nil {
		return nil, errors.New("not found")
	}
	uid, ok := r.byPersonID[personID]
	if !ok {
		return nil, errors.New("not found")
	}
	return r.Get(ctx, uid)
}

func (r *fakeUserRepo) SetZenodoTokens(context.Context, uuid.UUID, string, string) error {
	return errors.New("not implemented")
}

func (r *fakeUserRepo) ClearZenodoTokens(context.Context, uuid.UUID) error {
	return errors.New("not implemented")
}

// Fake person repo stores ent.Person to test mapping logic via transform
type fakePersonRepo struct {
	byID map[uuid.UUID]*ent.Person
}

func (r *fakePersonRepo) Get(_ context.Context, id uuid.UUID) (*entities.Person, error) {
	p, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return transform.ToEntityPtr[entities.Person](p), nil
}

func (r *fakePersonRepo) GetByEmail(_ context.Context, email string) (*entities.Person, error) {
	for _, p := range r.byID {
		if strings.EqualFold(p.Email, email) {
			return transform.ToEntityPtr[entities.Person](p), nil
		}
	}
	return nil, errors.New("not found")
}

func (r *fakePersonRepo) List(context.Context) ([]*entities.Person, error) {
	out := make([]*ent.Person, 0, len(r.byID))
	for _, p := range r.byID {
		out = append(out, p)
	}
	return transform.ToEntitiesPtr[entities.Person](out), nil
}

func (r *fakePersonRepo) Search(_ context.Context, _ string, limit int) ([]entities.Person, error) {
	out := make([]*ent.Person, 0)
	for _, p := range r.byID {
		out = append(out, p)
		if len(out) >= limit {
			break
		}
	}
	return transform.ToEntities[entities.Person](out), nil
}

func (r *fakePersonRepo) SetORCID(_ context.Context, personID uuid.UUID, orcidID string) error {
	p, ok := r.byID[personID]
	if !ok {
		return errors.New("not found")
	}
	p.OrcidID = orcidID
	return nil
}

func (r *fakePersonRepo) ClearORCID(_ context.Context, personID uuid.UUID) error {
	p, ok := r.byID[personID]
	if !ok {
		return errors.New("not found")
	}
	p.OrcidID = ""
	return nil
}

func (r *fakePersonRepo) Create(_ context.Context, p entities.Person) (*entities.Person, error) {
	if r.byID == nil {
		r.byID = map[uuid.UUID]*ent.Person{}
	}
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	ep := &ent.Person{
		ID:              p.ID,
		UserID:          p.UserID,
		Name:            p.Name,
		GivenName:       p.GivenName,
		FamilyName:      p.FamilyName,
		Email:           p.Email,
		AvatarURL:       p.AvatarUrl,
		Description:     p.Description,
		OrgCustomFields: p.OrgCustomFields,
	}
	if p.ORCiD != nil {
		ep.OrcidID = *p.ORCiD
	}
	r.byID[p.ID] = ep
	return transform.ToEntityPtr[entities.Person](ep), nil
}

func (r *fakePersonRepo) Update(_ context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error) {
	if r.byID == nil {
		r.byID = map[uuid.UUID]*ent.Person{}
	}
	ep := &ent.Person{
		ID:              id,
		UserID:          p.UserID,
		Name:            p.Name,
		GivenName:       p.GivenName,
		FamilyName:      p.FamilyName,
		Email:           p.Email,
		AvatarURL:       p.AvatarUrl,
		Description:     p.Description,
		OrgCustomFields: p.OrgCustomFields,
	}
	if p.ORCiD != nil {
		ep.OrcidID = *p.ORCiD
	}
	r.byID[id] = ep
	return transform.ToEntityPtr[entities.Person](ep), nil
}

func (r *fakePersonRepo) Delete(id uuid.UUID) error {
	if r.byID != nil {
		delete(r.byID, id)
	}
	return nil
}

func TestService_Search_DelegatesToClient(t *testing.T) {
	users := &fakeUserRepo{}
	people := &fakePersonRepo{}
	client := &fakeOrcidClient{
		searchOut: []ext.OrcidPerson{
			{FirstName: "John", LastName: "Doe", ORCID: "0000-0001-2345-6789"},
		},
	}

	svc := orcid.NewService(users, people, client)

	got, err := svc.Search(context.Background(), "John Doe")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ORCID != "0000-0001-2345-6789" {
		t.Fatalf("unexpected ORCID: %s", got[0].ORCID)
	}
}

func TestService_Link_AlreadyLinked(t *testing.T) {
	userID := uuid.New()
	personID := uuid.New()
	existing := "0000-0001-9999-8888"

	users := &fakeUserRepo{
		byID: map[uuid.UUID]entities.User{
			userID: {ID: userID, PersonID: personID},
		},
		byPersonID: map[uuid.UUID]uuid.UUID{personID: userID},
	}

	people := &fakePersonRepo{
		byID: map[uuid.UUID]*ent.Person{
			personID: {ID: personID, OrcidID: existing},
		},
	}

	client := &fakeOrcidClient{exchangeOut: "0000-0002-1111-2222"}

	svc := orcid.NewService(users, people, client)

	err := svc.Link(context.Background(), userID, "AUTH_CODE")
	if !errors.Is(err, orcid.ErrAlreadyLinked) {
		t.Fatalf("expected ErrAlreadyLinked, got %v", err)
	}
}

func TestService_Link_SetsORCIDOnPerson(t *testing.T) {
	userID := uuid.New()
	personID := uuid.New()

	users := &fakeUserRepo{
		byID: map[uuid.UUID]entities.User{
			userID: {ID: userID, PersonID: personID},
		},
		byPersonID: map[uuid.UUID]uuid.UUID{personID: userID},
	}

	people := &fakePersonRepo{
		byID: map[uuid.UUID]*ent.Person{
			personID: {ID: personID},
		},
	}

	client := &fakeOrcidClient{exchangeOut: "0000-0002-1111-2222"}

	svc := orcid.NewService(users, people, client)

	if err := svc.Link(context.Background(), userID, "AUTH_CODE"); err != nil {
		t.Fatalf("Link failed: %v", err)
	}

	p, _ := people.Get(context.Background(), personID)
	if p.ORCiD == nil || *p.ORCiD != "0000-0002-1111-2222" {
		t.Fatalf("expected ORCID set to 0000-0002-1111-2222, got %#v", p.ORCiD)
	}
}

func TestService_Unlink_ClearsORCIDOnPerson(t *testing.T) {
	userID := uuid.New()
	personID := uuid.New()
	existing := "0000-0002-1111-2222"

	users := &fakeUserRepo{
		byID: map[uuid.UUID]entities.User{
			userID: {ID: userID, PersonID: personID},
		},
		byPersonID: map[uuid.UUID]uuid.UUID{personID: userID},
	}

	people := &fakePersonRepo{
		byID: map[uuid.UUID]*ent.Person{
			personID: {ID: personID, OrcidID: existing},
		},
	}

	client := &fakeOrcidClient{}

	svc := orcid.NewService(users, people, client)

	if err := svc.Unlink(context.Background(), userID); err != nil {
		t.Fatalf("Unlink failed: %v", err)
	}

	p, _ := people.Get(context.Background(), personID)
	if p.ORCiD != nil && *p.ORCiD != "" {
		t.Fatalf("expected ORCID cleared, got %#v", *p.ORCiD)
	}
}
