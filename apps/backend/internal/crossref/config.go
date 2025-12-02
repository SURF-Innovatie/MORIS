package crossref

import (
	"fmt"
	"os"
)

// Config holds the configuration for Crossref API
type Config struct {
	BaseURL   string
	UserAgent string
	Mailto    string
}

// GetConfig returns the Crossref configuration from environment variables
func GetConfig() (*Config, error) {
	userAgent := os.Getenv("CROSSREF_USER_AGENT")
	mailto := os.Getenv("CROSSREF_MAILTO")

	if userAgent == "" {
		userAgent = "MORIS/1.0 (https://github.com/SURF-Innovatie/MORIS)"
	}

	if mailto == "" {
		return nil, fmt.Errorf("CROSSREF_MAILTO must be set")
	}

	baseURL := os.Getenv("CROSSREF_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.crossref.org"
	}

	return &Config{
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Mailto:    mailto,
	}, nil
}
