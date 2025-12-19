package ror

// ClientOptions holds configuration for the ROR client.
type ClientOptions struct {
	BaseUrl string
}

// ClientOption is a functional option for configuring the client.
type ClientOption func(*ClientOptions)

// DefaultClientOptions returns the default options.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseUrl: "https://api.ror.org/organizations",
	}
}

// WithBaseUrl sets the base URL for the ROR API.
func WithBaseUrl(url string) ClientOption {
	return func(o *ClientOptions) {
		o.BaseUrl = url
	}
}
