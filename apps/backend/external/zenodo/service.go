package zenodo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/google/uuid"
)

var (
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrMissingCode        = errors.New("missing_authorization_code")
	ErrAlreadyLinked      = errors.New("zenodo_already_linked")
	ErrNotLinked          = errors.New("zenodo_not_linked")
	ErrMissingBucketURL   = errors.New("missing_bucket_url")
	ErrMissingAccessToken = errors.New("missing_access_token")
)

// Service provides high-level operations for Zenodo integration
type Service interface {
	// OAuth Operations
	GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error)
	Link(ctx context.Context, userID uuid.UUID, code string) error
	Unlink(ctx context.Context, userID uuid.UUID) error
	IsLinked(ctx context.Context, userID uuid.UUID) (bool, error)

	// Deposition Operations
	CreateDeposition(ctx context.Context, userID uuid.UUID) (*Deposition, error)
	GetDeposition(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error)
	UpdateDeposition(ctx context.Context, userID uuid.UUID, depositionID int, metadata *DepositionMetadata) (*Deposition, error)
	DeleteDeposition(ctx context.Context, userID uuid.UUID, depositionID int) error
	ListDepositions(ctx context.Context, userID uuid.UUID) ([]Deposition, error)

	// File Operations
	UploadFile(ctx context.Context, userID uuid.UUID, depositionID int, filename string, data io.Reader) (*DepositionFile, error)

	// Actions
	Publish(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error)
	NewVersion(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error)
}

type service struct {
	client       *ent.Client
	userSvc      user.Service
	httpClient   *http.Client
	zenodoClient *Client
}

// NewService creates a new Zenodo service
func NewService(client *ent.Client, userSvc user.Service, httpClient *http.Client) Service {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	cfg, _ := GetConfig() // Config might not be available at startup
	var zenodoClient *Client
	if cfg != nil {
		zenodoClient = NewClient(httpClient, cfg)
	}

	return &service{
		client:       client,
		userSvc:      userSvc,
		httpClient:   httpClient,
		zenodoClient: zenodoClient,
	}
}

// getConfig returns the Zenodo configuration, initializing the client if needed
func (s *service) getConfig() (*Config, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	if s.zenodoClient == nil {
		s.zenodoClient = NewClient(s.httpClient, cfg)
	}
	return cfg, nil
}

// getAccessToken retrieves the Zenodo access token for a user
func (s *service) getAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	person, err := s.client.Person.Get(ctx, usr.PersonID)
	if err != nil {
		return "", fmt.Errorf("failed to get person: %w", err)
	}

	if person.ZenodoAccessToken == "" {
		return "", ErrNotLinked
	}

	return person.ZenodoAccessToken, nil
}

func (s *service) GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error) {
	// Verify user exists
	if _, err := s.client.User.Get(ctx, userID); err != nil {
		return "", ErrUnauthenticated
	}

	cfg, err := s.getConfig()
	if err != nil {
		return "", err
	}

	// Use userID as state for CSRF protection
	return cfg.GenerateAuthURL(userID.String()), nil
}

func (s *service) Link(ctx context.Context, userID uuid.UUID, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	cfg, err := s.getConfig()
	if err != nil {
		return err
	}

	tokens, err := cfg.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	person, err := s.client.Person.Get(ctx, usr.PersonID)
	if err != nil {
		return fmt.Errorf("failed to get person: %w", err)
	}

	if person.ZenodoAccessToken != "" {
		return ErrAlreadyLinked
	}

	// Update person with Zenodo tokens
	update := s.client.Person.UpdateOneID(usr.PersonID).
		SetZenodoAccessToken(tokens.AccessToken)

	if tokens.RefreshToken != "" {
		update = update.SetZenodoRefreshToken(tokens.RefreshToken)
	}

	if _, err := update.Save(ctx); err != nil {
		return fmt.Errorf("failed to link Zenodo: %w", err)
	}

	return nil
}

func (s *service) Unlink(ctx context.Context, userID uuid.UUID) error {
	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Clear Zenodo tokens from person
	if _, err := s.client.Person.
		UpdateOneID(usr.PersonID).
		ClearZenodoAccessToken().
		ClearZenodoRefreshToken().
		Save(ctx); err != nil {
		return fmt.Errorf("failed to unlink Zenodo: %w", err)
	}

	return nil
}

func (s *service) IsLinked(ctx context.Context, userID uuid.UUID) (bool, error) {
	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	person, err := s.client.Person.Get(ctx, usr.PersonID)
	if err != nil {
		return false, fmt.Errorf("failed to get person: %w", err)
	}

	return person.ZenodoAccessToken != "", nil
}

func (s *service) CreateDeposition(ctx context.Context, userID uuid.UUID) (*Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.CreateDeposition(ctx, accessToken)
}

func (s *service) GetDeposition(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.GetDeposition(ctx, accessToken, depositionID)
}

func (s *service) UpdateDeposition(ctx context.Context, userID uuid.UUID, depositionID int, metadata *DepositionMetadata) (*Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.UpdateDeposition(ctx, accessToken, depositionID, metadata)
}

func (s *service) DeleteDeposition(ctx context.Context, userID uuid.UUID, depositionID int) error {
	if _, err := s.getConfig(); err != nil {
		return err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return err
	}

	return s.zenodoClient.DeleteDeposition(ctx, accessToken, depositionID)
}

func (s *service) ListDepositions(ctx context.Context, userID uuid.UUID) ([]Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.ListDepositions(ctx, accessToken)
}

func (s *service) UploadFile(ctx context.Context, userID uuid.UUID, depositionID int, filename string, data io.Reader) (*DepositionFile, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	// First get the deposition to obtain the bucket URL
	deposition, err := s.zenodoClient.GetDeposition(ctx, accessToken, depositionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deposition: %w", err)
	}

	if deposition.Links == nil || deposition.Links.Bucket == "" {
		return nil, ErrMissingBucketURL
	}

	return s.zenodoClient.UploadFile(ctx, accessToken, deposition.Links.Bucket, filename, data)
}

func (s *service) Publish(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.Publish(ctx, accessToken, depositionID)
}

func (s *service) NewVersion(ctx context.Context, userID uuid.UUID, depositionID int) (*Deposition, error) {
	if _, err := s.getConfig(); err != nil {
		return nil, err
	}

	accessToken, err := s.getAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.zenodoClient.NewVersion(ctx, accessToken, depositionID)
}
