package env

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	Port       string
}

var Global Config

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist (e.g. in production),
		// but if it does and fails, we might want to know.
		// For now, we assume environment variables might be set directly.
	}

	Global = Config{
		AppEnv:     getEnv("APP_ENV", false, "dev"),
		DBHost:     getEnv("DB_HOST", true, ""),
		DBUser:     getEnv("DB_USER", true, ""),
		DBPassword: getEnv("DB_PASSWORD", true, ""),
		DBName:     getEnv("DB_NAME", true, ""),
		DBPort:     getEnv("DB_PORT", true, ""),
		JWTSecret:  getEnv("JWT_SECRET", true, ""),
		Port:       getEnv("PORT", true, ""),
	}
}

func getEnv(key string, required bool, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		if required {
			log.Panicf("Missing required environment variable: %s", key)
		}
		return fallback
	}
	return val
}

// IsDev returns true if the environment is development (or not set)
func IsDev() bool {
	return strings.EqualFold(Global.AppEnv, "dev") || strings.EqualFold(Global.AppEnv, "development")
}

// IsProd returns true if the environment is production
func IsProd() bool {
	return strings.EqualFold(Global.AppEnv, "prod") || strings.EqualFold(Global.AppEnv, "production")
}
