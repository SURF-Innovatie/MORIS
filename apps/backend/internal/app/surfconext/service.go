package surfconext

import (
	"context"
	"errors"
	"fmt"

	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

var (
	ErrMissingCode = errors.New("missing_authorization_code")
	ErrNoEmail     = errors.New("no_email_claim")
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	LoginWithCode(ctx context.Context, code string) (string, *entities.UserAccount, error)
}

type service struct {
	client  Client
	authSvc coreauth.Service
}

func NewService(client Client, authSvc coreauth.Service) Service {
	return &service{client: client, authSvc: authSvc}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	return s.client.AuthURL(ctx)
}

func (s *service) LoginWithCode(ctx context.Context, code string) (string, *entities.UserAccount, error) {
	if code == "" {
		return "", nil, ErrMissingCode
	}

	claims, err := s.client.ExchangeCode(ctx, code)
	if err != nil {
		return "", nil, err
	}
	if claims == nil || claims.Email == "" {
		return "", nil, ErrNoEmail
	}

	token, user, err := s.authSvc.LoginByEmail(ctx, claims.Email)
	if err != nil {
		return "", nil, fmt.Errorf("login by email: %w", err)
	}

	return token, user, nil
}
