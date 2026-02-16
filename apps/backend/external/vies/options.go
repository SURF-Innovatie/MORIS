package vies

// ClientOptions holds configuration for the VIES client.
type ClientOptions struct {
	BaseUrl string
}

// DefaultClientOptions returns the default client options.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseUrl: "https://ec.europa.eu/taxation_customs/vies/rest-api",
	}
}

// ClientOption is a function type for setting client options.
type ClientOption func(*ClientOptions)

// WithBaseUrl sets the base URL for the client.
func WithBaseUrl(url string) ClientOption {
	return func(opts *ClientOptions) {
		opts.BaseUrl = url
	}
}
