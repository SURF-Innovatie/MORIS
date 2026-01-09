package raid

// Options holds the configuration for the RAiD client.
type Options struct {
	BaseURL  string
	AuthURL  string
	Username string
	Password string
}

// DefaultOptions returns the default options for the RAiD client.
func DefaultOptions() Options {
	return Options{
		BaseURL: "https://api.demo.raid.org.au/",
		AuthURL: "https://auth.demo.raid.org.au/realms/RAiD/protocol/openid-connect/token",
	}
}
