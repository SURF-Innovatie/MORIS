package kvk

// Config holds the configuration for the KVK API client
type Config struct {
	// BaseURL is the base URL for the KVK API (e.g., https://api.kvk.nl/test/api)
	BaseURL string
	// APIKey is the API key for authentication
	APIKey string
}
