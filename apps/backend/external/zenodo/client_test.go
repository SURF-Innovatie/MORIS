package zenodo_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/SURF-Innovatie/MORIS/external/zenodo"
)

// RoundTripFunc implements http.RoundTripper for testing
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req), nil }

// NewTestClient returns an *http.Client with Transport replaced to avoid network calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{Transport: RoundTripFunc(fn)}
}

func testOpts() zenodo.Options {
	opts := zenodo.DefaultOptions(true) // sandbox
	opts.ClientID = "test-client-id"
	opts.ClientSecret = "test-client-secret"
	opts.RedirectURL = "http://localhost/callback"
	return opts
}

func TestClient_AuthURL(t *testing.T) {
	opts := testOpts()
	c := zenodo.NewClient(nil, opts)

	u, err := c.AuthURL("test-state")
	if err != nil {
		t.Fatalf("AuthURL failed: %v", err)
	}

	if !strings.Contains(u, "sandbox.zenodo.org") {
		t.Fatalf("expected sandbox host, got %s", u)
	}
	if !strings.Contains(u, "client_id=test-client-id") {
		t.Fatalf("expected client_id in url, got %s", u)
	}
	if !strings.Contains(u, "redirect_uri=") {
		t.Fatalf("expected redirect_uri in url, got %s", u)
	}
	if !strings.Contains(u, "state=test-state") {
		t.Fatalf("expected state in url, got %s", u)
	}
	if !strings.Contains(u, "scope=") {
		t.Fatalf("expected scope in url, got %s", u)
	}
}

func TestClient_ExchangeCode(t *testing.T) {
	opts := testOpts()

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Token endpoint
		if req.Method == http.MethodPost && req.URL.String() == opts.TokenURL {
			bodyBytes, _ := io.ReadAll(req.Body)
			body := string(bodyBytes)

			if !strings.Contains(body, "grant_type=authorization_code") {
				t.Fatalf("expected authorization_code, got body: %s", body)
			}
			if !strings.Contains(body, "code=AUTH_CODE") {
				t.Fatalf("expected code=AUTH_CODE, got body: %s", body)
			}

			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"access_token":"AT",
					"refresh_token":"RT",
					"token_type":"bearer",
					"scope":"deposit:write deposit:actions"
				}`)),
				Header: make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := zenodo.NewClient(httpClient, opts)

	tok, err := c.ExchangeCode(context.Background(), "AUTH_CODE")
	if err != nil {
		t.Fatalf("ExchangeCode failed: %v", err)
	}
	if tok.AccessToken != "AT" {
		t.Fatalf("expected access token AT, got %q", tok.AccessToken)
	}
	if tok.RefreshToken != "RT" {
		t.Fatalf("expected refresh token RT, got %q", tok.RefreshToken)
	}
}

func TestClient_CreateDeposition(t *testing.T) {
	opts := testOpts()

	mockDeposition := zenodo.Deposition{
		ID:    12345,
		State: zenodo.StateInProgress,
		Links: &zenodo.DepositionLinks{
			Bucket:  "https://sandbox.zenodo.org/api/files/test-bucket",
			Publish: "https://sandbox.zenodo.org/api/deposit/depositions/12345/actions/publish",
		},
		Metadata: &zenodo.DepositionMetadata{
			PrereserveDOI: &zenodo.PrereserveDOI{
				DOI:   "10.5072/zenodo.12345",
				RecID: 12345,
			},
		},
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Verify auth header
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"unauthorized"}`)),
				Header:     make(http.Header),
			}
		}

		// Create deposition endpoint
		if req.Method == http.MethodPost && req.URL.Path == "/api/deposit/depositions" {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := zenodo.NewClient(httpClient, opts)

	dep, err := c.CreateDeposition(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("CreateDeposition failed: %v", err)
	}
	if dep.ID != 12345 {
		t.Fatalf("expected ID 12345, got %d", dep.ID)
	}
	if dep.Links == nil || dep.Links.Bucket == "" {
		t.Fatalf("expected bucket URL to be set")
	}
	if dep.Metadata == nil || dep.Metadata.PrereserveDOI == nil || dep.Metadata.PrereserveDOI.DOI == "" {
		t.Fatalf("expected prereserved DOI to be set")
	}
}

func TestClient_Publish(t *testing.T) {
	opts := testOpts()

	mockDeposition := zenodo.Deposition{
		ID:        12345,
		State:     zenodo.StateDone,
		Submitted: true,
		DOI:       "10.5072/zenodo.12345",
		DOIURL:    "https://doi.org/10.5072/zenodo.12345",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"unauthorized"}`)),
				Header:     make(http.Header),
			}
		}

		if req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/actions/publish") {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 202,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := zenodo.NewClient(httpClient, opts)

	dep, err := c.Publish(context.Background(), "test-token", 12345)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if dep.DOI != "10.5072/zenodo.12345" {
		t.Fatalf("expected DOI, got %s", dep.DOI)
	}
	if !dep.Submitted {
		t.Fatalf("expected Submitted=true")
	}
}

func TestClient_UploadFile(t *testing.T) {
	opts := testOpts()

	mockFile := zenodo.DepositionFile{
		ID:       "file-id-123",
		Filename: "test-file.txt",
		Filesize: 1024,
		Checksum: "abc123",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"unauthorized"}`)),
				Header:     make(http.Header),
			}
		}

		// File upload endpoint: PUT to bucket URL
		if req.Method == http.MethodPut && strings.Contains(req.URL.String(), "test-bucket") {
			respBytes, _ := json.Marshal(mockFile)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := zenodo.NewClient(httpClient, opts)

	fileData := bytes.NewBufferString("test file content")
	f, err := c.UploadFile(context.Background(), "test-token", "https://sandbox.zenodo.org/api/files/test-bucket", "test-file.txt", fileData)
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}
	if f.Filename != "test-file.txt" {
		t.Fatalf("expected filename test-file.txt, got %s", f.Filename)
	}
	if f.ID == "" {
		t.Fatalf("expected file id to be set")
	}
}

func TestClient_UpdateDeposition(t *testing.T) {
	opts := testOpts()

	mockDeposition := zenodo.Deposition{
		ID:    12345,
		State: zenodo.StateInProgress,
		Metadata: &zenodo.DepositionMetadata{
			Title:       "Test Upload",
			UploadType:  zenodo.UploadTypeDataset,
			Description: "A test description",
			Creators: []zenodo.Creator{
				{Name: "Doe, John", Affiliation: "Test University"},
			},
			AccessRight: zenodo.AccessOpen,
		},
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"unauthorized"}`)),
				Header:     make(http.Header),
			}
		}

		// Update endpoint
		if req.Method == http.MethodPut && strings.Contains(req.URL.Path, "/deposit/depositions/12345") {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message":"not found"}`)),
			Header:     make(http.Header),
		}
	})

	c := zenodo.NewClient(httpClient, opts)

	metadata := &zenodo.DepositionMetadata{
		Title:       "Test Upload",
		UploadType:  zenodo.UploadTypeDataset,
		Description: "A test description",
		Creators: []zenodo.Creator{
			{Name: "Doe, John", Affiliation: "Test University"},
		},
		AccessRight: zenodo.AccessOpen,
	}

	dep, err := c.UpdateDeposition(context.Background(), "test-token", 12345, metadata)
	if err != nil {
		t.Fatalf("UpdateDeposition failed: %v", err)
	}
	if dep.Metadata == nil || dep.Metadata.Title != "Test Upload" {
		t.Fatalf("expected title 'Test Upload', got %#v", dep.Metadata)
	}
}
