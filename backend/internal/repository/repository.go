package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"url-shortener-backend/internal/models"

	// Drivers – the one needed will be selected at runtime.
	_ "modernc.org/sqlite"
)

var (
	ErrNotFound = errors.New("url mapping not found")
)

// Repository defines persistence operations for URLs.
type Repository interface {
	FindByOriginalURL(ctx context.Context, originalURL string) (*models.URL, error)
	Create(ctx context.Context, originalURL, shortCode string) (*models.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	GetRecentURLs(ctx context.Context, limit int) ([]*models.URL, error)
	Close() error
}

// NewRepository creates a repository based on the driver and DSN.
func NewRepository(driver string, dsn string) (Repository, error) {
	// Ensure the directory exists for SQLite file if needed.
	if driver == "sqlite" {
		dir := filepath.Dir(dsn)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create sqlite directory: %w", err)
			}
		}
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s database: %w", driver, err)
	}

	// Connection pool tuning.
	if driver == "sqlite" {
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
	}
	db.SetConnMaxLifetime(time.Hour)

	repo := &sqlRepository{db: db, driver: driver}
	if err := repo.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return repo, nil
}

type sqlRepository struct {
	db     *sql.DB
	driver string
}

func (r *sqlRepository) migrate() error {
	var schema string
	if r.driver == "sqlite" {
		schema = `
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_url TEXT NOT NULL UNIQUE,
			short_code TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);
		`
	} else { // mysql
		schema = `
		CREATE TABLE IF NOT EXISTS urls (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			original_url VARCHAR(2048) NOT NULL UNIQUE,
			short_code VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);
		`
	}
	_, err := r.db.Exec(schema)
	return err
}

func (r *sqlRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE original_url = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, originalURL)
	var u models.URL
	var ts string
	if err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = parseTime(ts)
	return &u, nil
}

func (r *sqlRepository) Create(ctx context.Context, originalURL, shortCode string) (*models.URL, error) {
	query := `INSERT INTO urls (original_url, short_code) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, originalURL, shortCode)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.URL{ID: id, OriginalURL: originalURL, ShortCode: shortCode, CreatedAt: time.Now().UTC()}, nil
}

func (r *sqlRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, shortCode)
	var u models.URL
	var ts string
	if err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = parseTime(ts)
	return &u, nil
}

func (r *sqlRepository) GetRecentURLs(ctx context.Context, limit int) ([]*models.URL, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT id, original_url, short_code, created_at FROM urls ORDER BY id DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.URL
	for rows.Next() {
		var u models.URL
		var ts string
		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &ts); err != nil {
			return nil, err
		}
		u.CreatedAt = parseTime(ts)
		result = append(result, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *sqlRepository) Close() error {
	return r.db.Close()
}

func parseTime(raw string) time.Time {
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02T15:04:05Z"}
	for _, l := range layouts {
		if t, err := time.Parse(l, raw); err == nil {
			return t
		}
	}
	return time.Now().UTC()
}
