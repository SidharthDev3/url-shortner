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

	_ "modernc.org/sqlite"
)

var (
	ErrNotFound = errors.New("url mapping not found")
)

// Repository defines data persistence operations for URLs.
type Repository interface {
	FindByOriginalURL(ctx context.Context, originalURL string) (*models.URL, error)
	Create(ctx context.Context, originalURL, shortCode string) (*models.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	GetRecentURLs(ctx context.Context, limit int) ([]*models.URL, error)
	Close() error
}

// SQLiteRepository implements Repository backed by SQLite.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates or opens an SQLite database and runs migrations.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Optimize connection pooling for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	repo := &SQLiteRepository{db: db}
	if err := repo.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL UNIQUE,
		short_code TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
	CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);
	`
	_, err := r.db.Exec(schema)
	return err
}

// FindByOriginalURL looks up an existing URL mapping (Duplicate check).
func (r *SQLiteRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE original_url = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, originalURL)

	var u models.URL
	var createdAtStr string
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	u.CreatedAt = parseTime(createdAtStr)
	return &u, nil
}

// Create inserts a new originalURL and shortCode mapping into SQLite.
func (r *SQLiteRepository) Create(ctx context.Context, originalURL, shortCode string) (*models.URL, error) {
	query := `INSERT INTO urls (original_url, short_code) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, originalURL, shortCode)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.URL{
		ID:          id,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// FindByShortCode finds the original URL given a short code.
func (r *SQLiteRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, shortCode)

	var u models.URL
	var createdAtStr string
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	u.CreatedAt = parseTime(createdAtStr)
	return &u, nil
}

// GetRecentURLs fetches recently shortened URLs for the dashboard.
func (r *SQLiteRepository) GetRecentURLs(ctx context.Context, limit int) ([]*models.URL, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `SELECT id, original_url, short_code, created_at FROM urls ORDER BY id DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*models.URL
	for rows.Next() {
		var u models.URL
		var createdAtStr string
		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAtStr); err != nil {
			return nil, err
		}
		u.CreatedAt = parseTime(createdAtStr)
		urls = append(urls, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// Close closes the database handle.
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func parseTime(raw string) time.Time {
	// Attempt multiple SQLite datetime formats
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, raw); err == nil {
			return t
		}
	}
	return time.Now().UTC()
}
