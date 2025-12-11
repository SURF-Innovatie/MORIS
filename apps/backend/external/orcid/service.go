package orcid

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrMissingCode     = errors.New("missing_authorization_code")
	ErrAlreadyLinked   = errors.New("orcid_already_linked")
)

type Service interface {
	// GetAuthURL returns an ORCID OAuth authorization URL for the given user, and the state parameter.
	GetAuthURL(ctx context.Context, userID uuid.UUID) (string, string, error)

	// Link links an ORCID ID to the given user using the ORCID auth code.
	Link(ctx context.Context, userID uuid.UUID, code string) error

	// Unlink unlinks an ORCID ID from the given user.
	Unlink(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	client   *ent.Client
	userSvc  user.Service
	provider *oidc.Provider
	oauth2   *oauth2.Config
}

func NewService(client *ent.Client, userSvc user.Service, provider *oidc.Provider, oauth2 *oauth2.Config) Service {
	return &service{
		client:   client,
		userSvc:  userSvc,
		provider: provider,
		oauth2:   oauth2,
	}
}

func (s *service) GetAuthURL(ctx context.Context, userID uuid.UUID) (string, string, error) {
	if _, err := s.client.User.Get(ctx, userID); err != nil {
		return "", "", ErrUnauthenticated
	}

	state, err := generateState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	return s.oauth2.AuthCodeURL(state), state, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *service) Link(ctx context.Context, userID uuid.UUID, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	token, err := s.oauth2.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	// Extract ORCID ID.
	orcidID, ok := token.Extra("orcid").(string)
	if !ok || orcidID == "" {
		return errors.New("failed to retrieve ORCID ID")
	}

	usr, err := s.userSvc.GetAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if usr.Person.ORCiD != nil && *usr.Person.ORCiD != "" {
		return ErrAlreadyLinked
	}

	// Update person with ORCID ID
	if _, err := s.client.Person.
		UpdateOneID(usr.Person.Id).
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
