package crossref_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/external/crossref"
)

func TestClient_GetWork_SetsHeadersAndParses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/works/") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("User-Agent") != "UA" {
			t.Fatalf("missing/invalid User-Agent: %q", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("mailto") != "m@example.org" {
			t.Fatalf("missing/invalid mailto: %q", r.Header.Get("mailto"))
		}

		w.Header().Set("X-Rate-Limit-Limit", "2")
		w.Header().Set("X-Rate-Limit-Interval", "1s")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"ok",
			"message-type":"work",
			"message-version":"1.0.0",
			"message":{"DOI":"10.1234/abc","title":["T"]}
		}`))
	}))
	defer ts.Close()

	cfg := &crossref.Config{BaseURL: ts.URL, UserAgent: "UA", Mailto: "m@example.org"}
	c := crossref.NewClientWithHTTP(cfg, &http.Client{Timeout: 2 * time.Second})

	work, err := c.GetWork(context.Background(), "10.1234/abc")
	if err != nil {
		t.Fatalf("GetWork error: %v", err)
	}
	if work.DOI != "10.1234/abc" {
		t.Fatalf("expected DOI 10.1234/abc, got %q", work.DOI)
	}
	if len(work.Title) != 1 || work.Title[0] != "T" {
		t.Fatalf("unexpected title: %#v", work.Title)
	}
}

func TestClient_GetWork_404_ReturnsErrNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cfg := &crossref.Config{BaseURL: ts.URL, UserAgent: "UA", Mailto: "m@example.org"}
	c := crossref.NewClient(cfg)

	_, err := c.GetWork(context.Background(), "10.1/x")
	if err == nil || !errors.Is(err, crossref.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
