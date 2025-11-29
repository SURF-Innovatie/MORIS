package orcid

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/google/uuid"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrMissingCode     = errors.New("missing_authorization_code")
	ErrAlreadyLinked   = errors.New("orcid_already_linked")
)

type Service interface {
	// GetAuthURL returns an ORCID OAuth authorization URL for the given user.
	GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error)

	// Link links an ORCID ID to the given user using the ORCID auth code.
	Link(ctx context.Context, userID uuid.UUID, code string) error

	// Unlink unlinks an ORCID ID from the given user.
	Unlink(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	client  *ent.Client
	userSvc user.Service
}

func NewService(client *ent.Client, userSvc user.Service) Service {
	return &service{client, userSvc}
}

func (s *service) GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error) {
	// Optional: verify user exists
	if _, err := s.client.User.Get(ctx, userID); err != nil {
		return "", ErrUnauthenticated
	}

	cfg, err := GetORCIDConfig()
	if err != nil {
		return "", err
	}
	return cfg.GenerateAuthURL(), nil
}

func (s *service) Link(ctx context.Context, userID uuid.UUID, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	cfg, err := GetORCIDConfig()
	if err != nil {
		return err
	}

	orcidID, err := cfg.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if usr.ORCiD != nil {
		return ErrAlreadyLinked
	}

	// Update person with ORCID ID
	if _, err := s.client.Person.
		UpdateOneID(usr.PersonID).
		SetOrcidID(orcidID).
		Save(ctx); err != nil {
		return fmt.Errorf("failed to link ORCID ID: %w", err)
	}

	return nil
}

func (s *service) Unlink(ctx context.Context, userID uuid.UUID) error {
	usr, err := s.userSvc.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Clear ORCID ID from person
	if _, err := s.client.Person.
		UpdateOneID(usr.PersonID).
		ClearOrcidID().
		Save(ctx); err != nil {
		return fmt.Errorf("failed to unlink ORCID ID: %w", err)
	}
	return nil
}
