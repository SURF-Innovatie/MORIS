package csvsource

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
)

type CSVSource struct {
	filePath string
}

func NewCSVSource(filePath string) *CSVSource {
	return &CSVSource{filePath: filePath}
}

func (s *CSVSource) Name() string {
	return "csv"
}

func (s *CSVSource) DisplayName() string {
	return "CSV Import"
}

func (s *CSVSource) SupportedTypes() []adapter.DataType {
	return []adapter.DataType{adapter.DataTypeUser}
}

func (s *CSVSource) InputInfo() adapter.InputInfo {
	return adapter.InputInfo{
		Type:        adapter.TransferTypeFile,
		Label:       "Upload CSV",
		Description: "Import users from a CSV file (columns: name, email, orcid)",
		MimeType:    "text/csv",
	}
}

func (s *CSVSource) Connect(ctx context.Context) error {
	if s.filePath == "" {
		return errors.New("file path is required")
	}
	return nil
}

func (s *CSVSource) FetchProjects(ctx context.Context, opts adapter.FetchOptions) (<-chan events.Event, <-chan error) {
	out := make(chan events.Event)
	errs := make(chan error, 1)
	close(out)
	errs <- errors.New("CSV source does not support project import")
	close(errs)
	return out, errs
}

func (s *CSVSource) FetchUsers(ctx context.Context, opts adapter.FetchOptions) (<-chan *adapter.UserContext, <-chan error) {
	out := make(chan *adapter.UserContext)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		f, err := os.Open(s.filePath)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		// Skip header
		if _, err := reader.Read(); err != nil {
			if err != io.EOF {
				errs <- err
			}
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				record, err := reader.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					errs <- err
					return
				}

				// Map CSV record to UserContext (simplified)
				if len(record) < 2 {
					continue
				}

				// In a real implementation, we would create entities here
				// For this example, we just show the structure
				out <- &adapter.UserContext{
					// Person: ...
					// User: ...
				}
			}
		}
	}()

	return out, errs
}

func (s *CSVSource) Close() error {
	return nil
}
