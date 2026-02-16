package vies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Client is the client for the VIES VAT API.
type Client struct {
	httpClient *http.Client
	options    ClientOptions
}

// NewClient creates a new VIES API client.
func NewClient(httpClient *http.Client, opts ...ClientOption) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Client{
		httpClient: httpClient,
		options:    options,
	}
}

// CheckVatNumber validates a VAT number against the VIES API.
// The vatNumber should include the country code prefix (e.g., "NL822655287B01").
func (c *Client) CheckVatNumber(ctx context.Context, vatNumber string) (*VatCheckResponse, error) {
	// Extract country code (first 2 characters) and number (rest)
	if len(vatNumber) < 3 {
		return nil, fmt.Errorf("invalid VAT number: too short")
	}

	countryCode := strings.ToUpper(vatNumber[:2])
	number := vatNumber[2:]

	return c.CheckVatNumberWithCountry(ctx, countryCode, number)
}

// CheckVatNumberWithCountry validates a VAT number against the VIES API with explicit country code.
func (c *Client) CheckVatNumberWithCountry(ctx context.Context, countryCode, vatNumber string) (*VatCheckResponse, error) {
	reqBody := VatCheckRequest{
		CountryCode: strings.ToUpper(countryCode),
		VatNumber:   vatNumber,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqUrl := fmt.Sprintf("%s/check-vat-number", c.options.BaseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result VatCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
