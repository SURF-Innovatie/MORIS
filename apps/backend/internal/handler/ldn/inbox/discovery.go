// Package inbox provides LDN Inbox discovery middleware.
package inbox

import (
	"net/http"
	"os"
)

// DiscoveryMiddleware adds the LDN inbox Link header to responses.
// Per W3C LDN spec: https://www.w3.org/TR/ldn/#discovery
func DiscoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originURL := os.Getenv("LDN_ORIGIN_URL")
		if originURL == "" {
			originURL = "http://localhost:8080"
		}
		inboxURL := originURL + "/api/ldn/inbox"

		// Add Link header for inbox discovery
		w.Header().Add("Link", `<`+inboxURL+`>; rel="http://www.w3.org/ns/ldp#inbox"`)

		next.ServeHTTP(w, r)
	})
}

// InboxDiscoveryHandler handles HEAD requests for inbox discovery.
// Per LDN spec, clients can discover inbox via HTTP HEAD or GET.
func InboxDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	originURL := os.Getenv("LDN_ORIGIN_URL")
	if originURL == "" {
		originURL = "http://localhost:8080"
	}
	inboxURL := originURL + "/api/ldn/inbox"

	w.Header().Set("Link", `<`+inboxURL+`>; rel="http://www.w3.org/ns/ldp#inbox"`)
	w.WriteHeader(http.StatusOK)
}
