package crossref_test

import (
	"context"
	"errors"
	"testing"

	ex "github.com/SURF-Innovatie/MORIS/external/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/app/crossref"
)

type fakeClient struct {
	work    *ex.Work
	workErr error
}

func (f *fakeClient) GetWork(context.Context, string) (*ex.Work, error) { return f.work, f.workErr }
func (f *fakeClient) GetWorks(context.Context, string, int) ([]ex.Work, error) {
	return nil, errors.New("not used")
}
func (f *fakeClient) GetJournal(context.Context, string) (*ex.Journal, error) {
	return nil, errors.New("not used")
}
func (f *fakeClient) GetJournals(context.Context, string, int) ([]ex.Journal, error) {
	return nil, errors.New("not used")
}

func TestService_GetWork_MapsNotFound(t *testing.T) {
	fc := &fakeClient{workErr: ex.ErrNotFound}
	svc := crossref.NewService(fc)

	_, err := svc.GetWork(context.Background(), "10/x")
	if !errors.Is(err, crossref.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_GetWork_PassesThroughSuccess(t *testing.T) {
	fc := &fakeClient{work: &ex.Work{DOI: "10.1/abc"}}
	svc := crossref.NewService(fc)

	w, err := svc.GetWork(context.Background(), "10.1/abc")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if w.DOI != "10.1/abc" {
		t.Fatalf("unexpected DOI: %q", w.DOI)
	}
}
