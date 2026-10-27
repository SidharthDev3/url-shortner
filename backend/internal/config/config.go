package config

import (
    "os"
    "strings"
)

// Config holds runtime configuration for the backend.
type Config struct {
    Port     string // server port, e.g. "8080"
    DBPath   string // fallback SQLite file path (used if driver is sqlite)
    BaseURL  string // base URL for generated short links
    Env      string // development / production
    DBDriver string // "sqlite" or "mysql"
    DBDSN    string // full DSN for the selected driver
}

// Load reads env vars (with defaults) and returns a Config.
func Load() *Config {
    port := envOr("PORT", "8080")
    dbPath := envOr("DB_PATH", "./data/urls.db")
    baseURL := envOr("BASE_URL", "http://localhost:"+port)
    baseURL = strings.TrimSuffix(baseURL, "/")
    env := envOr("ENV", "development")
    dbDriver := envOr("DB_DRIVER", "sqlite")
    dbDSN := os.Getenv("DB_DSN")
    if dbDSN == "" && dbDriver == "sqlite" {
        dbDSN = dbPath // SQLite DSN is just the file path
    }
    return &Config{
        Port:     port,
        DBPath:   dbPath,
        BaseURL:  baseURL,
        Env:      env,
        DBDriver: dbDriver,
        DBDSN:    dbDSN,
    }
}

// envOr returns the env var value or the fallback.
func envOr(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
