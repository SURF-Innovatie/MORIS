package zenodo

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	ext "github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/google/uuid"
)

type fakeUserRepo struct {
	byID map[uuid.UUID]identity.User
}

func (r *fakeUserRepo) Get(_ context.Context, id uuid.UUID) (*identity.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	uu := u
	return &uu, nil
}

func (r *fakeUserRepo) GetByPersonID(context.Context, uuid.UUID) (*identity.User, error) {
	return nil, errors.New("not implemented")
}
func (r *fakeUserRepo) Create(context.Context, identity.User) (*identity.User, error) {
	return nil, errors.New("not implemented")
}
func (r *fakeUserRepo) Update(context.Context, uuid.UUID, identity.User) (*identity.User, error) {
	return nil, errors.New("not implemented")
}
func (r *fakeUserRepo) Delete(context.Context, uuid.UUID) error {
	return errors.New("not implemented")
}
func (r *fakeUserRepo) ToggleActive(context.Context, uuid.UUID, bool) error {
	return errors.New("not implemented")
}
func (r *fakeUserRepo) ListUsers(context.Context, int, int) ([]identity.User, int, error) {
	return nil, 0, errors.New("not implemented")
}

func (r *fakeUserRepo) SetZenodoTokens(_ context.Context, userID uuid.UUID, access, refresh string) error {
	u, ok := r.byID[userID]
	if !ok {
		return errors.New("not found")
	}
	u.ZenodoAccessToken = &access
	if refresh != "" {
		u.ZenodoRefreshToken = &refresh
	}
	r.byID[userID] = u
	return nil
}

func (r *fakeUserRepo) ClearZenodoTokens(_ context.Context, userID uuid.UUID) error {
	u, ok := r.byID[userID]
	if !ok {
		return errors.New("not found")
	}
	u.ZenodoAccessToken = nil
	u.ZenodoRefreshToken = nil
	r.byID[userID] = u
	return nil
}

type fakeZenodoClient struct {
	authURL string

	exchangeTok *ext.TokenResponse
	exchangeErr error

	lastAccessToken string
	createCalled    bool
}

func (c *fakeZenodoClient) AuthURL(state string) (string, error) {
	if state == "" {
		return "", errors.New("missing state")
	}
	return c.authURL, nil
}

func (c *fakeZenodoClient) ExchangeCode(context.Context, string) (*ext.TokenResponse, error) {
	if c.exchangeErr != nil {
		return nil, c.exchangeErr
	}
	return c.exchangeTok, nil
}

func (c *fakeZenodoClient) RefreshToken(context.Context, string) (*ext.TokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeZenodoClient) CreateDeposition(_ context.Context, accessToken string) (*ext.Deposition, error) {
	c.lastAccessToken = accessToken
	c.createCalled = true
	return &ext.Deposition{ID: 1}, nil
}

func (c *fakeZenodoClient) GetDeposition(context.Context, string, int) (*ext.Deposition, error) {
	return nil, errors.New("not implemented")
}
func (c *fakeZenodoClient) UpdateDeposition(context.Context, string, int, *ext.DepositionMetadata) (*ext.Deposition, error) {
	return nil, errors.New("not implemented")
}
func (c *fakeZenodoClient) DeleteDeposition(context.Context, string, int) error {
	return errors.New("not implemented")
}
func (c *fakeZenodoClient) ListDepositions(context.Context, string) ([]ext.Deposition, error) {
	return nil, errors.New("not implemented")
}
func (c *fakeZenodoClient) UploadFile(context.Context, string, string, string, io.Reader) (*ext.DepositionFile, error) {
	return nil, errors.New("not implemented")
}
func (c *fakeZenodoClient) Publish(context.Context, string, int) (*ext.Deposition, error) {
	return nil, errors.New("not implemented")
}
func (c *fakeZenodoClient) NewVersion(context.Context, string, int) (*ext.Deposition, error) {
	return nil, errors.New("not implemented")
}

func TestService_GetAuthURL_UsesUserIDAsState(t *testing.T) {
	userID := uuid.New()

	users := &fakeUserRepo{byID: map[uuid.UUID]identity.User{
		userID: {ID: userID, ZenodoAccessToken: nil, ZenodoRefreshToken: nil},
	}}
	client := &fakeZenodoClient{authURL: "https://example/auth?x=y"}

	svc := NewService(users, client)

	u, err := svc.GetAuthURL(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetAuthURL failed: %v", err)
	}
	if !strings.Contains(u, "https://example/auth") {
		t.Fatalf("unexpected auth url %s", u)
	}
}

func TestService_Link_SetsTokens(t *testing.T) {
	userID := uuid.New()

	users := &fakeUserRepo{byID: map[uuid.UUID]identity.User{
		userID: {ID: userID, ZenodoAccessToken: nil, ZenodoRefreshToken: nil},
	}}

	client := &fakeZenodoClient{
		exchangeTok: &ext.TokenResponse{AccessToken: "AT", RefreshToken: "RT"},
	}

	svc := NewService(users, client)

	if err := svc.Link(context.Background(), userID, "CODE"); err != nil {
		t.Fatalf("Link failed: %v", err)
	}
}

func TestService_Link_AlreadyLinked(t *testing.T) {
	userID := uuid.New()
	existing := "EXISTING"

	users := &fakeUserRepo{byID: map[uuid.UUID]identity.User{
		userID: {ID: userID, ZenodoAccessToken: &existing},
	}}
	client := &fakeZenodoClient{
		exchangeTok: &ext.TokenResponse{AccessToken: "AT"},
	}

	svc := NewService(users, client)

	err := svc.Link(context.Background(), userID, "CODE")
	if !errors.Is(err, ErrAlreadyLinked) {
		t.Fatalf("expected ErrAlreadyLinked, got %v", err)
	}
}

func TestService_Unlink_ClearsTokens(t *testing.T) {
	userID := uuid.New()
	at := "AT"
	rt := "RT"

	users := &fakeUserRepo{byID: map[uuid.UUID]identity.User{
		userID: {ID: userID, ZenodoAccessToken: &at, ZenodoRefreshToken: &rt},
	}}
	client := &fakeZenodoClient{}

	svc := NewService(users, client)

	if err := svc.Unlink(context.Background(), userID); err != nil {
		t.Fatalf("Unlink failed: %v", err)
	}

	p, _ := users.Get(context.Background(), userID)
	if p.ZenodoAccessToken != nil || p.ZenodoRefreshToken != nil {
		t.Fatalf("expected tokens cleared, got access=%#v refresh=%#v", p.ZenodoAccessToken, p.ZenodoRefreshToken)
	}
}

func TestService_CreateDeposition_UsesStoredAccessToken(t *testing.T) {
	userID := uuid.New()
	at := "AT"

	users := &fakeUserRepo{byID: map[uuid.UUID]identity.User{
		userID: {ID: userID, ZenodoAccessToken: &at},
	}}
	client := &fakeZenodoClient{}

	svc := NewService(users, client)

	_, err := svc.CreateDeposition(context.Background(), userID)
	if err != nil {
		t.Fatalf("CreateDeposition failed: %v", err)
	}
	if !client.createCalled {
		t.Fatalf("expected client.CreateDeposition to be called")
	}
	if client.lastAccessToken != "AT" {
		t.Fatalf("expected access token AT, got %q", client.lastAccessToken)
	}
}
