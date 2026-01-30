package doi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

var (
	ErrNotFound = errors.New("doi_not_found")
)

type Service interface {
	Resolve(ctx context.Context, doi string) (*dto.Work, error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
			// Default CheckRedirect follows redirects (up to 10)
		},
	}
}

func (s *service) Resolve(ctx context.Context, doi string) (*dto.Work, error) {
	// Clean DOI
	doi = strings.TrimPrefix(doi, "https://doi.org/")

	url := fmt.Sprintf("https://doi.org/%s", doi)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Request multiple formats via content negotiation
	// We prioritize CSL JSON (Crossref style) as it's cleaner, then JSON-LD (Schema.org)
	req.Header.Set("Accept", "application/vnd.citationstyles.csl+json, application/ld+json, application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	// Parse based on content type
	if strings.Contains(contentType, "application/vnd.citationstyles.csl+json") || strings.Contains(contentType, "application/json") {
		return s.parseCSLJSON(resp.Body, doi)
	} else if strings.Contains(contentType, "application/ld+json") {
		return s.parseJSONLD(resp.Body, doi)
	}

	return nil, fmt.Errorf("unsupported content type: %s", contentType)
}

// CSL JSON (Crossref-like) structure
type cslItem struct {
	DOI       string `json:"DOI"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	GraphType string `json:"@type"` // For JSON-LD sometimes served as JSON
	Publisher string `json:"publisher"`
	Author    []struct {
		Given  string `json:"given"`
		Family string `json:"family"`
		Name   string `json:"name"` // Sometimes used in JSON-LD
	} `json:"author"`
	Issued struct {
		DateParts [][]interface{} `json:"date-parts"` // e.g. [[2021, 1, 15]]
	} `json:"issued"`
}

func (s *service) parseCSLJSON(r io.Reader, originalDOI string) (*dto.Work, error) {
	var item cslItem
	if err := json.NewDecoder(r).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode CSL JSON: %w", err)
	}

	// Basic mapping
	w := &dto.Work{
		DOI:       item.DOI,
		Title:     item.Title,
		Publisher: item.Publisher,
		Type:      mapCSLType(item.Type),
	}

	if w.DOI == "" {
		w.DOI = originalDOI
	}

	// Authors
	for _, a := range item.Author {
		name := ""
		if a.Given != "" && a.Family != "" {
			name = fmt.Sprintf("%s %s", a.Given, a.Family)
		} else if a.Name != "" {
			name = a.Name
		} else if a.Family != "" {
			name = a.Family
		}

		if name != "" {
			w.Authors = append(w.Authors, name)
		}
	}

	// Date
	if len(item.Issued.DateParts) > 0 && len(item.Issued.DateParts[0]) > 0 {
		// Just taking the year/first part for simplicity or full date if available
		parts := item.Issued.DateParts[0]
		if len(parts) >= 1 {
			// Convert float/mix to int safe-ish
			val := fmt.Sprintf("%v", parts[0])
			w.Date = val
		}
	}

	return w, nil
}

// JSON-LD (Schema.org) structure generic map since it varies
func (s *service) parseJSONLD(r io.Reader, originalDOI string) (*dto.Work, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON-LD: %w", err)
	}

	w := &dto.Work{
		DOI: originalDOI, // Often not explicitly in the root of JSON-LD in a simple way, assume matches
	}

	// Title
	if name, ok := data["name"].(string); ok {
		w.Title = name
	} else if headline, ok := data["headline"].(string); ok {
		w.Title = headline
	}

	// Type
	if t, ok := data["@type"].(string); ok {
		w.Type = mapSchemaType(t)
	} else if tList, ok := data["@type"].([]interface{}); ok && len(tList) > 0 {
		if tStr, ok := tList[0].(string); ok {
			w.Type = mapSchemaType(tStr)
		}
	} else {
		w.Type = entities.Other
	}

	// Publisher
	if pub, ok := data["publisher"].(map[string]interface{}); ok {
		if name, ok := pub["name"].(string); ok {
			w.Publisher = name
		}
	}

	return w, nil
}

func mapCSLType(t string) entities.ProductType {
	switch strings.ToLower(t) {
	case "dataset":
		return entities.Dataset
	case "software", "computerprogram":
		return entities.Software
	case "graphic", "image":
		return entities.Image
	case "article", "journal-article", "proceedings-article":
		return entities.Other // Or a specific type if MORIS supports it, otherwise Other is safe
	case "book", "book-chapter":
		return entities.Other
	default:
		return entities.Other
	}
}

func mapSchemaType(t string) entities.ProductType {
	switch strings.ToLower(t) {
	case "dataset":
		return entities.Dataset
	case "softwareapplication", "softwaresourcecode":
		return entities.Software
	case "imageobject":
		return entities.Image
	default:
		return entities.Other
	}
}
