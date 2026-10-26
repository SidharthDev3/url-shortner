package config

import (
	"os"
	"strings"
)

// Config holds runtime configuration settings for the backend application.
type Config struct {
	Port    string
	DBPath  string
	BaseURL string
	Env     string
}

// Load loads configuration from environment variables with fallback defaults.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/urls.db"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:    port,
		DBPath:  dbPath,
		BaseURL: baseURL,
		Env:     env,
	}
}
