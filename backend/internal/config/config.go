package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration for the backend.
type Config struct {
	Port        string // server port, e.g. "8080"
	DBPath      string // fallback SQLite file path (used if driver is sqlite)
	BaseURL     string // base URL for generated short links
	Env         string // development / production
	DBDriver    string // "sqlite" or "mysql"
	DBDSN       string // full DSN for the selected driver
	DBTLSCAFile string // optional PEM bundle used to verify a MySQL server certificate
}

// Load reads an optional local .env file and environment variables. Existing
// process environment variables take precedence over values in .env.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	port := envOr("PORT", "8080")
	dbPath := envOr("DB_PATH", "./data/urls.db")
	baseURL := strings.TrimSuffix(envOr("BASE_URL", "http://localhost:"+port), "/")
	dbDriver := strings.ToLower(strings.TrimSpace(envOr("DB_DRIVER", "sqlite")))
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" && dbDriver == "sqlite" {
		dbDSN = dbPath // SQLite DSN is just the file path.
	}

	return &Config{
		Port:        port,
		DBPath:      dbPath,
		BaseURL:     baseURL,
		Env:         envOr("ENV", "development"),
		DBDriver:    dbDriver,
		DBDSN:       dbDSN,
		DBTLSCAFile: os.Getenv("DB_TLS_CA_FILE"),
	}, nil
}

// envOr returns the env var value or the fallback.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
