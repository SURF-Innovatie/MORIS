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

	doi "github.com/SURF-Innovatie/MORIS/external/doi"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
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
		},
	}
}

func (s *service) Resolve(ctx context.Context, originalDoi string) (*dto.Work, error) {
	parsedDoi, err := doi.Parse(originalDoi)

	url := fmt.Sprintf("https://doi.org/%s", parsedDoi)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	if strings.Contains(contentType, "application/vnd.citationstyles.csl+json") || strings.Contains(contentType, "application/json") {
		return s.parseCSLJSON(resp.Body, parsedDoi)
	} else if strings.Contains(contentType, "application/ld+json") {
		return s.parseJSONLD(resp.Body, parsedDoi)
	}

	return nil, fmt.Errorf("unsupported content type: %s", contentType)
}

// CSL JSON (Crossref-like) structure
type cslItem struct {
	DOI       string `json:"DOI"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	GraphType string `json:"@type"`
	Publisher string `json:"publisher"`
	Author    []struct {
		Given  string `json:"given"`
		Family string `json:"family"`
		Name   string `json:"name"`
		ORCID  string `json:"ORCID"`
	} `json:"author"`
	Issued struct {
		DateParts [][]interface{} `json:"date-parts"`
	} `json:"issued"`
}

func (s *service) parseCSLJSON(r io.Reader, parsedDoi doi.DOI) (*dto.Work, error) {
	var item cslItem
	if err := json.NewDecoder(r).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode CSL JSON: %w", err)
	}

	w := &dto.Work{
		DOI:       parsedDoi,
		Title:     item.Title,
		Publisher: item.Publisher,
		Type:      mapCSLType(item.Type),
	}

	for _, a := range item.Author {
		wa := dto.WorkAuthor{
			Given:  strings.TrimSpace(a.Given),
			Family: strings.TrimSpace(a.Family),
			Name:   strings.TrimSpace(a.Name),
			ORCID:  strings.TrimSpace(a.ORCID),
		}

		// If Name missing, derive it
		if wa.Name == "" {
			switch {
			case wa.Given != "" && wa.Family != "":
				wa.Name = wa.Given + " " + wa.Family
			case wa.Family != "":
				wa.Name = wa.Family
			}
		}

		// Normalize ORCID if it is a URL
		wa.ORCID = strings.TrimPrefix(wa.ORCID, "https://orcid.org/")
		wa.ORCID = strings.TrimPrefix(wa.ORCID, "http://orcid.org/")

		if wa.Name != "" || wa.ORCID != "" {
			w.Authors = append(w.Authors, wa)
		}
	}

	if len(item.Issued.DateParts) > 0 && len(item.Issued.DateParts[0]) > 0 {
		parts := item.Issued.DateParts[0]
		if len(parts) >= 1 {
			w.Date = fmt.Sprintf("%v", parts[0])
		}
	}

	return w, nil
}

func (s *service) parseJSONLD(r io.Reader, parsedDoi doi.DOI) (*dto.Work, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON-LD: %w", err)
	}

	w := &dto.Work{
		DOI: parsedDoi,
	}

	if name, ok := data["name"].(string); ok {
		w.Title = name
	} else if headline, ok := data["headline"].(string); ok {
		w.Title = headline
	}

	if t, ok := data["@type"].(string); ok {
		w.Type = mapSchemaType(t)
	} else if tList, ok := data["@type"].([]interface{}); ok && len(tList) > 0 {
		if tStr, ok := tList[0].(string); ok {
			w.Type = mapSchemaType(tStr)
		}
	} else {
		w.Type = product.Other
	}

	if pub, ok := data["publisher"].(map[string]interface{}); ok {
		if name, ok := pub["name"].(string); ok {
			w.Publisher = name
		}
	}

	return w, nil
}

func mapCSLType(t string) product.ProductType {
	switch strings.ToLower(t) {
	case "dataset":
		return product.Dataset
	case "software", "computerprogram":
		return product.Software
	case "graphic", "image":
		return product.Image
	case "article", "journal-article", "proceedings-article":
		return product.Other
	case "book", "book-chapter":
		return product.Other
	default:
		return product.Other
	}
}

func mapSchemaType(t string) product.ProductType {
	switch strings.ToLower(t) {
	case "dataset":
		return product.Dataset
	case "softwareapplication", "softwaresourcecode":
		return product.Software
	case "imageobject":
		return product.Image
	default:
		return product.Other
	}
}
