package zenodo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/google/uuid"
)

var (
	ErrMissingCode   = errors.New("missing_authorization_code")
	ErrAlreadyLinked = errors.New("zenodo_already_linked")
	ErrNotLinked     = errors.New("zenodo_not_linked")
	ErrMissingBucket = errors.New("missing_bucket_url")
)

type Service interface {
	GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error)
	Link(ctx context.Context, userID uuid.UUID, code string) error
	Unlink(ctx context.Context, userID uuid.UUID) error
	IsLinked(ctx context.Context, userID uuid.UUID) (bool, error)

	CreateDeposition(ctx context.Context, userID uuid.UUID) (*zenodo.Deposition, error)
	GetDeposition(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error)
	UpdateDeposition(ctx context.Context, userID uuid.UUID, depositionID int, md *zenodo.DepositionMetadata) (*zenodo.Deposition, error)
	DeleteDeposition(ctx context.Context, userID uuid.UUID, depositionID int) error
	ListDepositions(ctx context.Context, userID uuid.UUID) ([]zenodo.Deposition, error)

	UploadFile(ctx context.Context, userID uuid.UUID, depositionID int, filename string, data io.Reader) (*zenodo.DepositionFile, error)

	Publish(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error)
	NewVersion(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error)
}

type Client interface {
	AuthURL(state string) (string, error)
	ExchangeCode(ctx context.Context, code string) (*zenodo.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*zenodo.TokenResponse, error)

	CreateDeposition(ctx context.Context, accessToken string) (*zenodo.Deposition, error)
	GetDeposition(ctx context.Context, accessToken string, depositionID int) (*zenodo.Deposition, error)
	UpdateDeposition(ctx context.Context, accessToken string, depositionID int, md *zenodo.DepositionMetadata) (*zenodo.Deposition, error)
	DeleteDeposition(ctx context.Context, accessToken string, depositionID int) error
	ListDepositions(ctx context.Context, accessToken string) ([]zenodo.Deposition, error)

	UploadFile(ctx context.Context, accessToken, bucketURL, filename string, data io.Reader) (*zenodo.DepositionFile, error)

	Publish(ctx context.Context, accessToken string, depositionID int) (*zenodo.Deposition, error)
	NewVersion(ctx context.Context, accessToken string, depositionID int) (*zenodo.Deposition, error)
}

type service struct {
	users  user.Repository
	client Client
}

func NewService(users user.Repository, client Client) Service {
	return &service{users: users, client: client}
}

func (s *service) GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error) {
	// ensure user exists (authz is handled by handlers/middleware)
	if _, err := s.users.Get(ctx, userID); err != nil {
		return "", err
	}
	return s.client.AuthURL(userID.String())
}

func (s *service) Link(ctx context.Context, userID uuid.UUID, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return err
	}

	if u.ZenodoAccessToken != nil && *u.ZenodoAccessToken != "" {
		return ErrAlreadyLinked
	}

	tok, err := s.client.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	if err := s.users.SetZenodoTokens(ctx, userID, tok.AccessToken, tok.RefreshToken); err != nil {
		return fmt.Errorf("set zenodo tokens: %w", err)
	}
	return nil
}

func (s *service) Unlink(ctx context.Context, userID uuid.UUID) error {
	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return err
	}
	return s.users.ClearZenodoTokens(ctx, u.ID)
}

func (s *service) IsLinked(ctx context.Context, userID uuid.UUID) (bool, error) {
	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return false, err
	}
	return u.ZenodoAccessToken != nil && *u.ZenodoAccessToken != "", nil
}

func (s *service) accessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return "", err
	}
	if u.ZenodoAccessToken == nil || *u.ZenodoAccessToken == "" {
		return "", ErrNotLinked
	}
	return *u.ZenodoAccessToken, nil
}

func (s *service) CreateDeposition(ctx context.Context, userID uuid.UUID) (*zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.CreateDeposition(ctx, tok)
}

func (s *service) GetDeposition(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.GetDeposition(ctx, tok, depositionID)
}

func (s *service) UpdateDeposition(ctx context.Context, userID uuid.UUID, depositionID int, md *zenodo.DepositionMetadata) (*zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.UpdateDeposition(ctx, tok, depositionID, md)
}

func (s *service) DeleteDeposition(ctx context.Context, userID uuid.UUID, depositionID int) error {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return err
	}
	return s.client.DeleteDeposition(ctx, tok, depositionID)
}

func (s *service) ListDepositions(ctx context.Context, userID uuid.UUID) ([]zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.ListDepositions(ctx, tok)
}

func (s *service) UploadFile(ctx context.Context, userID uuid.UUID, depositionID int, filename string, data io.Reader) (*zenodo.DepositionFile, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	dep, err := s.client.GetDeposition(ctx, tok, depositionID)
	if err != nil {
		return nil, err
	}
	if dep.Links == nil || dep.Links.Bucket == "" {
		return nil, ErrMissingBucket
	}

	return s.client.UploadFile(ctx, tok, dep.Links.Bucket, filename, data)
}

func (s *service) Publish(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.Publish(ctx, tok, depositionID)
}

func (s *service) NewVersion(ctx context.Context, userID uuid.UUID, depositionID int) (*zenodo.Deposition, error) {
	tok, err := s.accessToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.client.NewVersion(ctx, tok, depositionID)
}
