package orcid

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid network calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestSearch(t *testing.T) {
	// Setup environment for config
	os.Setenv("ORCID_CLIENT_ID", "test-client-id")
	os.Setenv("ORCID_CLIENT_SECRET", "test-client-secret")
	os.Setenv("ORCID_REDIRECT_URL", "http://localhost/callback")
	os.Setenv("ORCID_SANDBOX", "true")
	defer func() {
		os.Unsetenv("ORCID_CLIENT_ID")
		os.Unsetenv("ORCID_CLIENT_SECRET")
		os.Unsetenv("ORCID_REDIRECT_URL")
		os.Unsetenv("ORCID_SANDBOX")
	}()

	// Mock responses
	mockDetails := []PersonExpandedSearchResult{
		{
			OrcidID:     "0000-0001-2345-6789",
			GivenNames:  "John",
			FamilyNames: "Doe",
			CreditName:  "J. Doe",
		},
	}
	mockSearchResponse := expandedSearchResponse{
		ExpandedResult: mockDetails,
	}

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Mock Token Request
		if req.URL.Path == "/oauth/token" && req.Method == "POST" {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"access_token": "mock-token",
					"token_type": "bearer",
					"expires_in": 631138518,
					"scope": "/read-public",
					"orcid": null
				}`)),
				Header: make(http.Header),
			}
		}

		// Mock Search Request
		if strings.Contains(req.URL.Path, "/expanded-search") && req.Method == "GET" {
			// Verify Auth Header
			auth := req.Header.Get("Authorization")
			if auth != "Bearer mock-token" {
				return &http.Response{
					StatusCode: 401,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "unauthorized"}`)),
				}
			}

			respBytes, _ := json.Marshal(mockSearchResponse)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
		}
	})

	// Initialize service with nil ent client and user svc as they are not used in Search
	svc := NewService(nil, nil, client)

	// Execute Search
	results, err := svc.Search(context.Background(), "John Doe")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Verify Results
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].FirstName != "John" {
		t.Errorf("Expected FirstName John, got %s", results[0].FirstName)
	}
	if results[0].ORCID != "0000-0001-2345-6789" {
		t.Errorf("Expected ORCID 0000-0001-2345-6789, got %s", results[0].ORCID)
	}
}
