package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/SURF-Innovatie/MORIS/internal/env"
	"github.com/SURF-Innovatie/MORIS/internal/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5/middleware"
)

// ErrorLoggingMiddleware returns a middleware that logs errors and panics
func ErrorLoggingMiddleware(svc errorlog.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Container to capture explicit errors passed from handlers
			errContainer := &httputil.ErrorDetailsContainer{}
			// Inject container into context
			ctx := context.WithValue(r.Context(), httputil.ContextKeyErrorDetails, errContainer)
			r = r.WithContext(ctx)

			// Wrap ResponseWriter to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				// Handle Panic
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())
					msg := fmt.Sprintf("PANIC: %v", rec)

					// Log the panic
					userId := httputil.GetUserIDFromContext(r.Context())
					svc.Log(context.Background(), userId, r.Method, r.URL.Path, http.StatusInternalServerError, msg, stack)

					// detailed error for dev, generic for prod
					respMsg := msg
					if env.IsProd() {
						respMsg = "Internal Server Error"
					}

					// Attempt to write a 500 Internal Server Error response if not already written

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					// Re-use standard structure
					resp := httputil.BackendError{
						Code:    http.StatusInternalServerError,
						Status:  "Internal Server Error",
						Message: respMsg,
					}
					// Only show errors/stack in dev
					if env.IsDev() {
						resp.Errors = stack
					}

					// Use standard encoder
					_ = httputil.WriteJSON(w, http.StatusInternalServerError, resp)
					return
				}

				// Handle Standard Errors (non-panic)
				// Check if we have an error status code or explicit error details
				if ww.Status() >= 400 {
					// Check if we captured details from httputil.WriteError
					msg := errContainer.Message
					details := errContainer.Errors

					// If no explicit details were captured, but status is error, use status text
					if msg == "" {
						msg = http.StatusText(ww.Status())
					}

					// Log it
					userId := httputil.GetUserIDFromContext(r.Context())

					// Convert details to string if possible for storage
					var detailStr string
					if details != nil {
						detailStr = fmt.Sprintf("%v", details)
					}

					svc.Log(context.Background(), userId, r.Method, r.URL.Path, ww.Status(), msg, detailStr)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
