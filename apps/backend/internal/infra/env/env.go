package env

import (
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppEnv         string   `env:"APP_ENV" env-default:"dev"`
	DBHost         string   `env:"DB_HOST" env-required:"true"`
	DBUser         string   `env:"DB_USER" env-required:"true"`
	DBPassword     string   `env:"DB_PASSWORD" env-required:"true"`
	DBName         string   `env:"DB_NAME" env-required:"true"`
	DBPort         string   `env:"DB_PORT" env-required:"true"`
	CacheHost      string   `env:"CACHE_HOST" env-required:"true"`
	CachePort      string   `env:"CACHE_PORT" env-required:"true"`
	CachePassword  string   `env:"CACHE_PASSWORD"`
	CacheUser      string   `env:"CACHE_USER"`
	JWTSecret      string   `env:"JWT_SECRET" env-required:"true"`
	Port           string   `env:"PORT" env-required:"true"`
	RAiDAPIKey     string   `env:"RAID_API_KEY" env-default:"test-key"`
	APIBasePath    string   `env:"API_BASE_PATH" env-default:"/api"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" env-default:"*"`

	ORCID      ORCIDConfig      `env-prefix:"ORCID_"`
	Surfconext SurfconextConfig `env-prefix:"SURFCONEXT_"`
	Zenodo     ZenodoConfig     `env-prefix:"ZENODO_"`
	Crossref   CrossrefConfig   `env-prefix:"CROSSREF_"`
}

type ORCIDConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL"`
	Sandbox      bool   `env:"SANDBOX" env-default:"false"`
}

type SurfconextConfig struct {
	IssuerURL    string `env:"ISSUER_URL"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL"`
	Scopes       string `env:"SCOPES"`
}

type ZenodoConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL"`
	Sandbox      bool   `env:"SANDBOX" env-default:"false"`
}

type CrossrefConfig struct {
	BaseURL   string `env:"BASE_URL" env-default:"https://api.crossref.org"`
	UserAgent string `env:"USER_AGENT" env-default:"MORIS/1.0 (https://github.com/SURF-Innovatie/MORIS)"`
	Mailto    string `env:"MAILTO" env-required:"true"`
}

var Global Config

func init() {
	// Check if .env file exists
	if _, err := os.Stat(".env"); err == nil {
		// Load from .env file (and override with env vars)
		if err := cleanenv.ReadConfig(".env", &Global); err != nil {
			log.Fatal().Err(err).Msg("failed to load environment variables from .env")
		}
	} else {
		// Load from environment variables only
		if err := cleanenv.ReadEnv(&Global); err != nil {
			log.Fatal().Err(err).Msg("failed to load environment variables")
		}
	}
}

// IsDev returns true if the environment is development (or not set)
func IsDev() bool {
	return strings.EqualFold(Global.AppEnv, "dev") || strings.EqualFold(Global.AppEnv, "development")
}

// IsProd returns true if the environment is production
func IsProd() bool {
	return strings.EqualFold(Global.AppEnv, "prod") || strings.EqualFold(Global.AppEnv, "production")
}
