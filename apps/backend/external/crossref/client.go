package crossref

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrNotFound = errors.New("crossref_not_found")
)

type Client interface {
	GetWork(ctx context.Context, doi string) (*Work, error)
	GetWorks(ctx context.Context, query string, limit int) ([]Work, error)
	GetJournal(ctx context.Context, issn string) (*Journal, error)
	GetJournals(ctx context.Context, query string, limit int) ([]Journal, error)
}

type client struct {
	httpClient        *http.Client
	config            *Config
	lastRequest       time.Time
	rateLimitLimit    int
	rateLimitInterval int
}

func NewClient(config *Config) Client {
	return &client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		config:     config,

		// safe defaults
		rateLimitLimit:    1,
		rateLimitInterval: 1,
	}
}

func NewClientWithHTTP(config *Config, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &client{
		httpClient:        httpClient,
		config:            config,
		rateLimitLimit:    1,
		rateLimitInterval: 1,
	}
}

func (c *client) GetWork(ctx context.Context, doi string) (*Work, error) {
	u := fmt.Sprintf("%s/works/%s", c.config.BaseURL, url.PathEscape(doi))

	if err := c.maybeDelay(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	c.processResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		var out WorkResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode WorkResponse: %w", err)
		}
		return &out.Message, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *client) GetWorks(ctx context.Context, query string, limit int) ([]Work, error) {
	if limit <= 0 {
		limit = 20
	}

	u := fmt.Sprintf("%s/works?query=%s&rows=%d",
		c.config.BaseURL,
		url.QueryEscape(query),
		limit,
	)

	if err := c.maybeDelay(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	c.processResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		var out MultipleWorksResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode MultipleWorksResponse: %w", err)
		}
		return out.Message.Items, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *client) GetJournal(ctx context.Context, issn string) (*Journal, error) {
	u := fmt.Sprintf("%s/journals/%s", c.config.BaseURL, url.PathEscape(issn))

	if err := c.maybeDelay(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	c.processResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		var out JournalResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode JournalResponse: %w", err)
		}
		return &out.Message, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *client) GetJournals(ctx context.Context, query string, limit int) ([]Journal, error) {
	if limit <= 0 {
		limit = 20
	}

	u := fmt.Sprintf("%s/journals?query=%s&rows=%d",
		c.config.BaseURL,
		url.QueryEscape(query),
		limit,
	)

	if err := c.maybeDelay(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	c.processResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		var out MultipleJournalsResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode MultipleJournalsResponse: %w", err)
		}
		return out.Message.Items, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *client) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.config.UserAgent)
	// If you want to keep this behavior:
	req.Header.Set("mailto", c.config.Mailto)
}

func (c *client) processResponse(resp *http.Response) {
	if limitStr := resp.Header.Get("X-Rate-Limit-Limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			c.rateLimitLimit = limit
		}
	}
	if c.rateLimitLimit <= 0 {
		c.rateLimitLimit = 1
	}

	if intervalStr := resp.Header.Get("X-Rate-Limit-Interval"); intervalStr != "" {
		if len(intervalStr) > 0 && intervalStr[len(intervalStr)-1] == 's' {
			intervalStr = intervalStr[:len(intervalStr)-1]
		}
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			c.rateLimitInterval = interval
		}
	}
	if c.rateLimitInterval <= 0 {
		c.rateLimitInterval = 1
	}

	c.lastRequest = time.Now()
}

func (c *client) maybeDelay() error {
	if c.lastRequest.IsZero() {
		return nil
	}
	delayMs := float64(c.rateLimitInterval) / float64(c.rateLimitLimit) * 1000
	delay := time.Duration(math.Ceil(delayMs)) * time.Millisecond
	time.Sleep(delay)
	return nil
}
