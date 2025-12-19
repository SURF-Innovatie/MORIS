package doi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// ModernDoiRegex matches standard DOI formats: 10.xxxx/yyyy
	ModernDoiRegex = regexp.MustCompile(`(?i)^10\.\d{4,9}/[-._;()/:a-z0-9<>]+$`)
	// OldDoiRegex matches older Wiley DOI formats: 10.1002/xxxx
	OldDoiRegex = regexp.MustCompile(`(?i)^10\.1002/[^\s]+$`)
)

// DOI represents a Digital Object Identifier.
// It is always stored in its normalized form (without https://doi.org/ prefix).
type DOI struct {
	value string
}

// Parse attempts to parse a string into a DOI object.
// It normalizes the input and validates it against standard DOI patterns.
func Parse(s string) (DOI, error) {
	normalized, err := Normalize(s)
	if err != nil {
		return DOI{}, err
	}

	if !IsValid(normalized) {
		return DOI{}, errors.New("invalid DOI format")
	}

	return DOI{value: normalized}, nil
}

// Normalize takes a DOI string (potentially a URL) and returns the raw DOI.
// Example: "https://doi.org/10.1000/182" -> "10.1000/182"
func Normalize(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("empty DOI string")
	}

	// If it doesn't start with http/https, assume it's already a raw DOI (or try to be)
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return strings.Trim(s, "/"), nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if !strings.EqualFold(u.Host, "doi.org") && !strings.EqualFold(u.Host, "dx.doi.org") {
		return "", errors.New("URL is not a doi.org link")
	}

	// Path usually starts with /, trim it
	return strings.Trim(u.Path, "/"), nil
}

// IsValid checks if the provided string matches known DOI patterns.
// The input should be a normalized DOI string, not a URL.
func IsValid(s string) bool {
	return ModernDoiRegex.MatchString(s) || OldDoiRegex.MatchString(s)
}

// String returns the string representation of the DOI.
func (d DOI) String() string {
	return d.value
}

// IsZero returns true if the DOI is empty (uninitialized).
func (d DOI) IsZero() bool {
	return d.value == ""
}

// MarshalJSON marshals the DOI to a JSON string.
func (d DOI) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.value)
}

// UnmarshalJSON unmarshals a JSON string into a DOI object.
// It performs validation during unmarshaling.
func (d *DOI) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}
