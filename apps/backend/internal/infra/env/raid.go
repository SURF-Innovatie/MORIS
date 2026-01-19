package env

import (
	"os"

	"github.com/SURF-Innovatie/MORIS/external/raid"
)

func RaidOptionsFromEnv() raid.Options {
	baseURL := os.Getenv("RAID_API_BASE_URL")
	authURL := os.Getenv("RAID_API_AUTH_URL")
	username := os.Getenv("RAID_USERNAME")
	password := os.Getenv("RAID_PASSWORD")

	opts := raid.DefaultOptions()
	if baseURL != "" {
		opts.BaseURL = baseURL
	}
	if authURL != "" {
		opts.AuthURL = authURL
	}
	opts.Username = username
	opts.Password = password

	return opts
}
