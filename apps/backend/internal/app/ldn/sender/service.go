// Package sender provides the LDN sender service for sending notifications to external inboxes.
// Implements W3C LDN Sender: https://www.w3.org/TR/ldn/#sender
package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/ldn"
)

// Service defines the interface for sending LDN notifications.
type Service interface {
	// Send sends an AS2 Activity to a target's LDN Inbox.
	Send(ctx context.Context, activity *ldn.Activity) error

	// DiscoverInbox discovers an LDN Inbox for a given resource/service URL.
	DiscoverInbox(ctx context.Context, resourceURL string) (string, error)
}

// sender implements the LDN sender service.
type sender struct {
	client    *http.Client
	originURL string
}

// NewService creates a new LDN sender service.
func NewService() Service {
	originURL := os.Getenv("LDN_ORIGIN_URL")
	if originURL == "" {
		originURL = "http://localhost:8080"
	}

	return &sender{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		originURL: originURL,
	}
}

// Send sends an AS2 Activity to the target's LDN Inbox.
func (s *sender) Send(ctx context.Context, activity *ldn.Activity) error {
	if activity.Target == nil || activity.Target.Inbox == "" {
		return fmt.Errorf("target inbox URL is required")
	}

	// Ensure origin is set
	if activity.Origin == nil {
		activity.Origin = ldn.NewService(s.originURL)
	}

	// Serialize activity to JSON-LD
	payload, err := json.Marshal(activity)
	if err != nil {
		return fmt.Errorf("failed to serialize activity: %w", err)
	}

	// Create POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, activity.Target.Inbox, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers per LDN spec
	req.Header.Set("Content-Type", "application/ld+json")
	req.Header.Set("Accept", "application/ld+json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	// Check for success (2xx status codes)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("inbox returned error status: %d", resp.StatusCode)
	}

	return nil
}

// DiscoverInbox discovers an LDN Inbox for a given resource/service URL.
// Per LDN spec, inbox is advertised via Link header or RDF.
func (s *sender) DiscoverInbox(ctx context.Context, resourceURL string) (string, error) {
	// Try HEAD request first (more efficient)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, resourceURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to discover inbox: %w", err)
	}
	defer resp.Body.Close()

	// Look for Link header with ldp:inbox relation
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		inbox := parseLinkHeader(linkHeader)
		if inbox != "" {
			return inbox, nil
		}
	}

	// If HEAD didn't provide inbox, try GET
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, resourceURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err = s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to discover inbox: %w", err)
	}
	defer resp.Body.Close()

	linkHeader = resp.Header.Get("Link")
	if linkHeader != "" {
		inbox := parseLinkHeader(linkHeader)
		if inbox != "" {
			return inbox, nil
		}
	}

	return "", fmt.Errorf("no inbox found for resource: %s", resourceURL)
}

// parseLinkHeader extracts the inbox URL from a Link header.
// Format: <URL>; rel="http://www.w3.org/ns/ldp#inbox"
func parseLinkHeader(header string) string {
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "ldp#inbox") {
			// Extract URL between < and >
			start := strings.Index(part, "<")
			end := strings.Index(part, ">")
			if start >= 0 && end > start {
				return part[start+1 : end]
			}
		}
	}
	return ""
}
