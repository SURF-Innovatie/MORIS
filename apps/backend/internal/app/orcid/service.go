package orcid

import (
	"context"
	"errors"

	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/google/uuid"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrMissingCode     = errors.New("missing_authorization_code")
	ErrAlreadyLinked   = errors.New("orcid_already_linked")
)

type Service interface {
	GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error)
	Link(ctx context.Context, userID uuid.UUID, code string) error
	Unlink(ctx context.Context, userID uuid.UUID) error
	Search(ctx context.Context, query string) ([]exorcid.OrcidPerson, error)
}

type service struct {
	users  user.Repository
	people person.Repository
	client OrcidClient
}

func NewService(users user.Repository, people person.Repository, client OrcidClient) Service {
	return &service{users: users, people: people, client: client}
}

func (s *service) GetAuthURL(ctx context.Context, userID uuid.UUID) (string, error) {
	// verify user exists
	if _, err := s.users.Get(ctx, userID); err != nil {
		return "", ErrUnauthenticated
	}
	return s.client.AuthURL()
}

func (s *service) Link(ctx context.Context, userID uuid.UUID, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return ErrUnauthenticated
	}

	p, err := s.people.Get(ctx, u.PersonID)
	if err != nil {
		return err
	}

	if p.ORCiD != nil && *p.ORCiD != "" {
		return ErrAlreadyLinked
	}

	orcidID, err := s.client.ExchangeCode(ctx, code)
	if err != nil {
		return err
	}

	return s.people.SetORCID(ctx, p.ID, orcidID)
}

func (s *service) Unlink(ctx context.Context, userID uuid.UUID) error {
	u, err := s.users.Get(ctx, userID)
	if err != nil {
		return ErrUnauthenticated
	}
	return s.people.ClearORCID(ctx, u.PersonID)
}

func (s *service) Search(ctx context.Context, query string) ([]exorcid.OrcidPerson, error) {
	return s.client.SearchExpanded(ctx, query)
}
