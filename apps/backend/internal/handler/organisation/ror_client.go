package organisation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
)

type RORItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// ROR returns country/addresses in nested structures, simplifying for now
	Addresses []struct {
		City string `json:"city"`
	} `json:"addresses"`
	Country struct {
		CountryName string `json:"country_name"`
	} `json:"country"`
}

type RORResponse struct {
	Items []RORItem `json:"id"` // ROR API actually returns different structure depending on version?
	// ROR V2: "items": [...]
	// Let's check ROR API V2 response.
	// https://api.ror.org/organizations?query=...
	// Returns: { "number_of_results": 10, "items": [...] }
}

// Internal structs for parsing ROR API V2
type rorApiResponse struct {
	Items []rorItemRaw `json:"items"`
}

type rorItemRaw struct {
	ID    string `json:"id"`
	Names []struct {
		Value string   `json:"value"`
		Types []string `json:"types"`
	} `json:"names"`
	Locations []struct {
		GeonamesDetails struct {
			CountryName string `json:"country_name"`
			Name        string `json:"name"`
		} `json:"geonames_details"`
	} `json:"locations"`
}

func SearchROR(query string) ([]RORItem, error) {
	if query == "" {
		return []RORItem{}, nil
	}

	// ROR API URL
	apiURL := fmt.Sprintf("https://api.ror.org/organizations?query=%s", url.QueryEscape(query))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ror api returned status: %d", resp.StatusCode)
	}

	var result rorApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Map to simplified RORItem for frontend
	items := make([]RORItem, 0, len(result.Items))
	for _, raw := range result.Items {
		name := ""
		// legitimate name selection strategy (prefer ror_display)
		for _, n := range raw.Names {
			if containsResult(n.Types, "ror_display") {
				name = n.Value
				break
			}
		}
		if name == "" && len(raw.Names) > 0 {
			name = raw.Names[0].Value
		}

		item := RORItem{
			ID:   raw.ID,
			Name: name,
		}

		// Map location to addresses/country structure expected by frontend
		if len(raw.Locations) > 0 {
			loc := raw.Locations[0].GeonamesDetails
			item.Country.CountryName = loc.CountryName
			item.Addresses = []struct {
				City string `json:"city"`
			}{{City: loc.Name}}
		}

		items = append(items, item)
	}

	return items, nil
}

func containsResult(slice []string, val string) bool {
	return slices.Contains(slice, val)
}
