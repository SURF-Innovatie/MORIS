package crossref

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Service defines the interface for Crossref API operations
type Service interface {
	// GetWork retrieves a single work by DOI
	GetWork(ctx context.Context, doi string) (*Work, error)

	// GetWorks searches for works by query string
	GetWorks(ctx context.Context, query string, limit int) ([]Work, error)

	// GetJournal retrieves a single journal by ISSN
	GetJournal(ctx context.Context, issn string) (*Journal, error)

	// GetJournals searches for journals by query string
	GetJournals(ctx context.Context, query string, limit int) ([]Journal, error)
}

type service struct {
	httpClient        *http.Client
	config            *Config
	lastRequest       time.Time
	rateLimitLimit    int
	rateLimitInterval int
}

// NewService creates a new Crossref service
func NewService(config *Config) Service {
	return &service{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config:            config,
		rateLimitLimit:    1,
		rateLimitInterval: 1,
	}
}

// NewServiceWithClient creates a new Crossref service with a custom HTTP client
func NewServiceWithClient(config *Config, httpClient *http.Client) Service {
	return &service{
		httpClient:        httpClient,
		config:            config,
		rateLimitLimit:    1,
		rateLimitInterval: 1,
	}
}

// GetWork retrieves a single work by DOI
func (s *service) GetWork(ctx context.Context, doi string) (*Work, error) {
	urlPath := fmt.Sprintf("%s/works/%s", s.config.BaseURL, url.PathEscape(doi))

	if !s.lastRequest.IsZero() {
		if err := s.doDelay(); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("mailto", s.config.Mailto)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	s.processResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result WorkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to deserialize WorkResponse for DOI %s: %v", doi, err)
		return nil, fmt.Errorf("failed to deserialize WorkResponse: %w", err)
	}

	return &result.Message, nil
}

// GetWorks searches for works by query string
func (s *service) GetWorks(ctx context.Context, query string, limit int) ([]Work, error) {
	if limit <= 0 {
		limit = 20
	}

	urlPath := fmt.Sprintf("%s/works?query=%s&rows=%d",
		s.config.BaseURL,
		url.QueryEscape(query),
		limit,
	)

	if !s.lastRequest.IsZero() {
		if err := s.doDelay(); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("mailto", s.config.Mailto)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	s.processResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result MultipleWorksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to deserialize MultipleWorksResponse for query %s: %v", query, err)
		return nil, fmt.Errorf("failed to deserialize MultipleWorksResponse: %w", err)
	}

	return result.Message.Items, nil
}

// GetJournal retrieves a single journal by ISSN
func (s *service) GetJournal(ctx context.Context, issn string) (*Journal, error) {
	urlPath := fmt.Sprintf("%s/journals/%s", s.config.BaseURL, url.PathEscape(issn))

	if !s.lastRequest.IsZero() {
		if err := s.doDelay(); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("mailto", s.config.Mailto)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	s.processResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result JournalResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to deserialize JournalResponse for ISSN %s: %v", issn, err)
		return nil, fmt.Errorf("failed to deserialize JournalResponse: %w", err)
	}

	return &result.Message, nil
}

// GetJournals searches for journals by query string
func (s *service) GetJournals(ctx context.Context, query string, limit int) ([]Journal, error) {
	if limit <= 0 {
		limit = 20
	}

	urlPath := fmt.Sprintf("%s/journals?query=%s&rows=%d",
		s.config.BaseURL,
		url.QueryEscape(query),
		limit,
	)

	if !s.lastRequest.IsZero() {
		if err := s.doDelay(); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("mailto", s.config.Mailto)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	s.processResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result MultipleJournalsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to deserialize MultipleJournalsResponse for query %s: %v", query, err)
		return nil, fmt.Errorf("failed to deserialize MultipleJournalsResponse: %w", err)
	}

	return result.Message.Items, nil
}

// processResponse extracts rate limit information from response headers
func (s *service) processResponse(resp *http.Response) {
	if limitStr := resp.Header.Get("X-Rate-Limit-Limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			s.rateLimitLimit = limit
		}
	}
	if s.rateLimitLimit == 0 {
		s.rateLimitLimit = 1
	}

	if intervalStr := resp.Header.Get("X-Rate-Limit-Interval"); intervalStr != "" {
		// Remove trailing 's' if present
		if len(intervalStr) > 0 && intervalStr[len(intervalStr)-1] == 's' {
			intervalStr = intervalStr[:len(intervalStr)-1]
		}
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			s.rateLimitInterval = interval
		}
	}
	if s.rateLimitInterval == 0 {
		s.rateLimitInterval = 1
	}

	s.lastRequest = time.Now()
}

// doDelay implements rate limiting delay
func (s *service) doDelay() error {
	delayMs := float64(s.rateLimitInterval) / float64(s.rateLimitLimit) * 1000
	delayDuration := time.Duration(math.Ceil(delayMs)) * time.Millisecond
	time.Sleep(delayDuration)
	return nil
}
