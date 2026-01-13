package orcid_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/SURF-Innovatie/MORIS/external/orcid"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req), nil }

func newTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{Transport: RoundTripFunc(fn)}
}

func TestClient_SearchExpanded_UsesClientCredentialsAndBearerToken(t *testing.T) {
	opts := orcid.Options{
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost/callback",
		BaseURL:       "https://sandbox.orcid.org",
		PublicBaseURL: "https://pub.sandbox.orcid.org/v3.0",
	}

	// Prepare JSON body matching ORCID expanded-search response schema.
	mockBody := map[string]any{
		"expanded-result": []map[string]any{
			{
				"orcid-id":     "0000-0001-2345-6789",
				"given-names":  "John",
				"family-names": "Doe",
				"credit-name":  "J. Doe",
			},
		},
	}
	mockRespBytes, _ := json.Marshal(mockBody)

	httpClient := newTestHTTPClient(func(req *http.Request) *http.Response {
		// Token request (client credentials)
		if req.Method == http.MethodPost && req.URL.String() == "https://sandbox.orcid.org/oauth/token" {
			b, _ := io.ReadAll(req.Body)
			body := string(b)

			if !strings.Contains(body, "grant_type=client_credentials") {
				t.Fatalf("expected grant_type=client_credentials, got body: %s", body)
			}
			if !strings.Contains(body, "scope=%2Fread-public") && !strings.Contains(body, "scope=/read-public") {
				t.Fatalf("expected scope /read-public, got body: %s", body)
			}

			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"access_token":"mock-token",
					"token_type":"bearer",
					"expires_in":3600,
					"scope":"/read-public"
				}`)),
				Header: make(http.Header),
			}
		}

		// Search request
		if req.Method == http.MethodGet && strings.HasPrefix(req.URL.String(), "https://pub.sandbox.orcid.org/v3.0/expanded-search") {
			if got := req.Header.Get("Authorization"); got != "Bearer mock-token" {
				t.Fatalf("expected Authorization 'Bearer mock-token', got %q", got)
			}
			if got := req.Header.Get("Accept"); got != "application/vnd.orcid+json" {
				t.Fatalf("expected Accept application/vnd.orcid+json, got %q", got)
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(mockRespBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := orcid.NewClient(httpClient, opts)

	results, err := c.SearchExpanded(context.Background(), "John Doe")
	if err != nil {
		t.Fatalf("SearchExpanded failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ORCID != "0000-0001-2345-6789" {
		t.Fatalf("expected ORCID 0000-0001-2345-6789, got %s", results[0].ORCID)
	}
	if results[0].FirstName != "John" {
		t.Fatalf("expected FirstName John, got %s", results[0].FirstName)
	}
	if results[0].LastName != "Doe" {
		t.Fatalf("expected LastName Doe, got %s", results[0].LastName)
	}
}

func TestClient_ExchangeCode_ReturnsORCID(t *testing.T) {
	opts := orcid.Options{
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost/callback",
		BaseURL:       "https://sandbox.orcid.org",
		PublicBaseURL: "https://pub.sandbox.orcid.org/v3.0",
	}

	httpClient := newTestHTTPClient(func(req *http.Request) *http.Response {
		if req.Method == http.MethodPost && req.URL.String() == "https://sandbox.orcid.org/oauth/token" {
			b, _ := io.ReadAll(req.Body)
			body := string(b)

			if !strings.Contains(body, "grant_type=authorization_code") {
				t.Fatalf("expected grant_type=authorization_code, got body: %s", body)
			}
			if !strings.Contains(body, "code=AUTH_CODE") {
				t.Fatalf("expected code=AUTH_CODE, got body: %s", body)
			}
			if !strings.Contains(body, "redirect_uri=http%3A%2F%2Flocalhost%2Fcallback") &&
				!strings.Contains(body, "redirect_uri=http://localhost/callback") {
				t.Fatalf("expected redirect_uri, got body: %s", body)
			}

			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"orcid":"0000-0002-1111-2222",
					"name":"Test User"
				}`)),
				Header: make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := orcid.NewClient(httpClient, opts)

	orcidID, err := c.ExchangeCode(context.Background(), "AUTH_CODE")
	if err != nil {
		t.Fatalf("ExchangeCode failed: %v", err)
	}
	if orcidID != "0000-0002-1111-2222" {
		t.Fatalf("expected ORCID 0000-0002-1111-2222, got %s", orcidID)
	}
}

func TestClient_AuthURL_BuildsAuthorizeURL(t *testing.T) {
	opts := orcid.Options{
		ClientID:      "test-client-id",
		RedirectURL:   "http://localhost/callback",
		BaseURL:       "https://sandbox.orcid.org",
		PublicBaseURL: "https://pub.sandbox.orcid.org/v3.0",
	}

	c := orcid.NewClient(nil, opts)

	u, err := c.AuthURL()
	if err != nil {
		t.Fatalf("AuthURL returned error: %v", err)
	}

	if !strings.HasPrefix(u, "https://sandbox.orcid.org/oauth/authorize?") {
		t.Fatalf("unexpected prefix: %s", u)
	}
	if !strings.Contains(u, "client_id=test-client-id") {
		t.Fatalf("expected client_id in url: %s", u)
	}
	if !strings.Contains(u, "redirect_uri=http%3A%2F%2Flocalhost%2Fcallback") &&
		!strings.Contains(u, "redirect_uri=http://localhost/callback") {
		t.Fatalf("expected redirect_uri in url: %s", u)
	}
	if !strings.Contains(u, "scope=%2Fauthenticate") && !strings.Contains(u, "scope=/authenticate") {
		t.Fatalf("expected scope /authenticate in url: %s", u)
	}
	if !strings.Contains(u, "response_type=code") {
		t.Fatalf("expected response_type=code in url: %s", u)
	}
}
