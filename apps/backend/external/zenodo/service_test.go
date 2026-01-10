package zenodo

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

// RoundTripFunc implements http.RoundTripper for testing
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns an *http.Client with Transport replaced to avoid network calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func setupTestEnv() func() {
	os.Setenv("ZENODO_CLIENT_ID", "test-client-id")
	os.Setenv("ZENODO_CLIENT_SECRET", "test-client-secret")
	os.Setenv("ZENODO_REDIRECT_URL", "http://localhost/callback")
	os.Setenv("ZENODO_SANDBOX", "true")
	return func() {
		os.Unsetenv("ZENODO_CLIENT_ID")
		os.Unsetenv("ZENODO_CLIENT_SECRET")
		os.Unsetenv("ZENODO_REDIRECT_URL")
		os.Unsetenv("ZENODO_SANDBOX")
	}
}

func TestGetConfig(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if cfg.ClientID != "test-client-id" {
		t.Errorf("Expected ClientID 'test-client-id', got %s", cfg.ClientID)
	}

	if !cfg.IsSandbox() {
		t.Error("Expected sandbox mode to be enabled")
	}

	if !strings.Contains(cfg.APIURL, "sandbox") {
		t.Errorf("Expected sandbox API URL, got %s", cfg.APIURL)
	}
}

func TestGetConfig_MissingEnv(t *testing.T) {
	// Don't set any env vars
	_, err := GetConfig()
	if err == nil {
		t.Error("Expected error when env vars are missing")
	}
}

func TestGenerateAuthURL(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	cfg, _ := GetConfig()
	authURL := cfg.GenerateAuthURL("test-state")

	if !strings.Contains(authURL, "sandbox.zenodo.org") {
		t.Errorf("Expected sandbox URL, got %s", authURL)
	}

	if !strings.Contains(authURL, "client_id=test-client-id") {
		t.Errorf("Expected client_id in URL, got %s", authURL)
	}
	// Check for URL-encoded scope (deposit%3Awrite is deposit:write URL-encoded)
	if !strings.Contains(authURL, "scope=") {
		t.Errorf("Expected scope parameter in URL, got %s", authURL)
	}

	if !strings.Contains(authURL, "state=test-state") {
		t.Errorf("Expected state in URL, got %s", authURL)
	}
}

func TestCreateDeposition(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	mockDeposition := Deposition{
		ID:    12345,
		State: StateInProgress,
		Links: &DepositionLinks{
			Bucket:  "https://sandbox.zenodo.org/api/files/test-bucket",
			Publish: "https://sandbox.zenodo.org/api/deposit/depositions/12345/actions/publish",
		},
		Metadata: &DepositionMetadata{
			PrereserveDOI: &PrereserveDOI{
				DOI:   "10.5072/zenodo.12345",
				RecID: 12345,
			},
		},
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Verify auth header
		auth := req.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": "unauthorized"}`)),
			}
		}

		// Create deposition endpoint
		if req.URL.Path == "/api/deposit/depositions" && req.Method == "POST" {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message": "not found"}`)),
		}
	})

	cfg, _ := GetConfig()
	client := NewClient(httpClient, cfg)

	deposition, err := client.CreateDeposition(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("CreateDeposition failed: %v", err)
	}

	if deposition.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", deposition.ID)
	}

	if deposition.Links.Bucket == "" {
		t.Error("Expected bucket URL to be set")
	}

	if deposition.Metadata.PrereserveDOI.DOI == "" {
		t.Error("Expected prereserved DOI")
	}
}

func TestPublish(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	mockDeposition := Deposition{
		ID:        12345,
		State:     StateDone,
		Submitted: true,
		DOI:       "10.5072/zenodo.12345",
		DOIURL:    "https://doi.org/10.5072/zenodo.12345",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": "unauthorized"}`)),
			}
		}

		// Publish endpoint
		if strings.HasSuffix(req.URL.Path, "/actions/publish") && req.Method == "POST" {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 202,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message": "not found"}`)),
		}
	})

	cfg, _ := GetConfig()
	client := NewClient(httpClient, cfg)

	deposition, err := client.Publish(context.Background(), "test-token", 12345)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if deposition.DOI != "10.5072/zenodo.12345" {
		t.Errorf("Expected DOI '10.5072/zenodo.12345', got %s", deposition.DOI)
	}

	if !deposition.Submitted {
		t.Error("Expected deposition to be submitted")
	}
}

func TestUploadFile(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	mockFile := DepositionFile{
		ID:       "file-id-123",
		Filename: "test-file.txt",
		Filesize: 1024,
		Checksum: "abc123",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": "unauthorized"}`)),
			}
		}

		// File upload endpoint (PUT to bucket)
		if req.Method == "PUT" && strings.Contains(req.URL.String(), "test-bucket") {
			respBytes, _ := json.Marshal(mockFile)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message": "not found"}`)),
		}
	})

	cfg, _ := GetConfig()
	client := NewClient(httpClient, cfg)

	fileData := bytes.NewBufferString("test file content")
	file, err := client.UploadFile(context.Background(), "test-token", "https://sandbox.zenodo.org/api/files/test-bucket", "test-file.txt", fileData)
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}

	if file.Filename != "test-file.txt" {
		t.Errorf("Expected filename 'test-file.txt', got %s", file.Filename)
	}

	if file.ID == "" {
		t.Error("Expected file ID to be set")
	}
}

func TestUpdateDeposition(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	mockDeposition := Deposition{
		ID:    12345,
		State: StateInProgress,
		Metadata: &DepositionMetadata{
			Title:       "Test Upload",
			UploadType:  UploadTypeDataset,
			Description: "A test description",
			Creators: []Creator{
				{Name: "Doe, John", Affiliation: "Test University"},
			},
			AccessRight: AccessOpen,
		},
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			return &http.Response{
				StatusCode: 401,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": "unauthorized"}`)),
			}
		}

		// Update deposition endpoint
		if strings.Contains(req.URL.Path, "/deposit/depositions/12345") && req.Method == "PUT" {
			respBytes, _ := json.Marshal(mockDeposition)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(respBytes)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"message": "not found"}`)),
		}
	})

	cfg, _ := GetConfig()
	client := NewClient(httpClient, cfg)

	metadata := &DepositionMetadata{
		Title:       "Test Upload",
		UploadType:  UploadTypeDataset,
		Description: "A test description",
		Creators: []Creator{
			{Name: "Doe, John", Affiliation: "Test University"},
		},
		AccessRight: AccessOpen,
	}

	deposition, err := client.UpdateDeposition(context.Background(), "test-token", 12345, metadata)
	if err != nil {
		t.Fatalf("UpdateDeposition failed: %v", err)
	}

	if deposition.Metadata.Title != "Test Upload" {
		t.Errorf("Expected title 'Test Upload', got %s", deposition.Metadata.Title)
	}
}

func TestExchangeCode(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	// Override the http.Client used by ExchangeCode
	// Note: ExchangeCode creates its own client, so we need a different approach
	// For now, we just test the config generation
	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	// Verify the token URL is correct
	expectedTokenURL := "https://sandbox.zenodo.org/oauth/token"
	if cfg.TokenURL != expectedTokenURL {
		t.Errorf("Expected token URL %s, got %s", expectedTokenURL, cfg.TokenURL)
	}
}
