package env

import (
	"os"
	"strings"
)

// IsDev returns true if the environment is development (or not set)
func IsDev() bool {
	env := os.Getenv("APP_ENV")
	// Default to dev if not set, or if explicitly set to dev/development
	return env == "" || strings.EqualFold(env, "dev") || strings.EqualFold(env, "development")
}

// IsProd returns true if the environment is production
func IsProd() bool {
	env := os.Getenv("APP_ENV")
	return strings.EqualFold(env, "prod") || strings.EqualFold(env, "production")
}
