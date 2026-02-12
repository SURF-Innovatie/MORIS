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

	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
)

var (
	ErrNotFound = errors.New("doi_not_found")
)

type Client interface {
	Resolve(ctx context.Context, doi string) (*Work, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

type ClientOption func(*client)

func WithHTTPClient(c *http.Client) ClientOption {
	return func(cl *client) {
		if c != nil {
			cl.httpClient = c
		}
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(cl *client) {
		if strings.TrimSpace(baseURL) != "" {
			cl.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://doi.org",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type WorkAuthor struct {
	Given  string
	Family string
	Name   string
	ORCID  string
}

type Work struct {
	DOI       DOI
	Title     string
	Type      product.ProductType
	Date      string
	Publisher string
	Authors   []WorkAuthor
}

func (c *client) Resolve(ctx context.Context, originalDoi string) (*Work, error) {
	parsedDoi, err := Parse(originalDoi)
	if err != nil {
		return nil, fmt.Errorf("failed to parse doi: %w", err)
	}

	url := fmt.Sprintf("%s/%s", c.baseURL, parsedDoi)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.citationstyles.csl+json, application/ld+json, application/json")

	resp, err := c.httpClient.Do(req)
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
	switch {
	case strings.Contains(contentType, "application/vnd.citationstyles.csl+json"),
		strings.Contains(contentType, "application/json"):
		return parseCSLJSON(resp.Body, parsedDoi)

	case strings.Contains(contentType, "application/ld+json"):
		return parseJSONLD(resp.Body, parsedDoi)

	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

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

func parseCSLJSON(r io.Reader, parsedDoi DOI) (*Work, error) {
	var item cslItem
	if err := json.NewDecoder(r).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode CSL JSON: %w", err)
	}

	w := &Work{
		DOI:       parsedDoi,
		Title:     item.Title,
		Publisher: item.Publisher,
		Type:      mapCSLType(item.Type),
	}

	for _, a := range item.Author {
		wa := WorkAuthor{
			Given:  strings.TrimSpace(a.Given),
			Family: strings.TrimSpace(a.Family),
			Name:   strings.TrimSpace(a.Name),
			ORCID:  strings.TrimSpace(a.ORCID),
		}

		if wa.Name == "" {
			switch {
			case wa.Given != "" && wa.Family != "":
				wa.Name = wa.Given + " " + wa.Family
			case wa.Family != "":
				wa.Name = wa.Family
			}
		}

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

func parseJSONLD(r io.Reader, parsedDoi DOI) (*Work, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON-LD: %w", err)
	}

	w := &Work{
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
