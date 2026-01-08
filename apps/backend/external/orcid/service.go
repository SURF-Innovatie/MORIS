package orcid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
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

	// Search searches for people in the ORCID public registry.
	Search(ctx context.Context, query string) ([]OrcidPerson, error)
}

type service struct {
	client     *ent.Client
	userSvc    user.Service
	httpClient *http.Client
}

func NewService(client *ent.Client, userSvc user.Service, httpClient *http.Client) Service {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &service{client, userSvc, httpClient}
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

	usr, err := s.userSvc.GetAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if usr.Person.ORCiD != nil && *usr.Person.ORCiD != "" {
		return ErrAlreadyLinked
	}

	// Update person with ORCID ID
	if _, err := s.client.Person.
		UpdateOneID(usr.Person.ID).
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

func (s *service) Search(ctx context.Context, query string) ([]OrcidPerson, error) {
	cfg, err := GetORCIDConfig()
	if err != nil {
		return nil, err
	}

	token, err := cfg.GetClientCredentialsToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	baseURL := cfg.GetPublicAPIURL()
	// Use expanded-search for better performance/details
	searchURL := fmt.Sprintf("%s/expanded-search?q=%s", baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.orcid+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var response expandedSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	results := make([]OrcidPerson, len(response.ExpandedResult))
	for i, r := range response.ExpandedResult {
		results[i] = r.ToPerson()
	}

	return results, nil
}
