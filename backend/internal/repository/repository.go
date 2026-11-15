package repository

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"url-shortener-backend/internal/models"

	// Drivers – the one needed will be selected at runtime.
	mysql "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFound = errors.New("url mapping not found")
)

// RegisterMySQLTLSConfig registers a certificate authority bundle for a MySQL
// DSN that uses tls=<name>, for example tls=rds.
func RegisterMySQLTLSConfig(name, caFile string) error {
	pemData, err := os.ReadFile(caFile)
	if err != nil {
		return fmt.Errorf("read MySQL CA bundle: %w", err)
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pemData) {
		return fmt.Errorf("no certificates found in MySQL CA bundle %q", caFile)
	}

	return mysql.RegisterTLSConfig(name, &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    roots,
	})
}

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
	driver = strings.ToLower(strings.TrimSpace(driver))
	if driver != "sqlite" && driver != "mysql" {
		return nil, fmt.Errorf("unsupported database driver %q", driver)
	}
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is required for %s", driver)
	}

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
			original_url TEXT NOT NULL,
			original_url_hash CHAR(64) CHARACTER SET ascii NOT NULL UNIQUE,
			short_code CHAR(6) CHARACTER SET ascii NOT NULL UNIQUE,
			created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
		);
		`
	}
	_, err := r.db.Exec(schema)
	return err
}

func (r *sqlRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE original_url = ? LIMIT 1`
	args := []any{originalURL}
	if r.driver == "mysql" {
		query = `SELECT id, original_url, short_code, created_at FROM urls WHERE original_url_hash = ? AND original_url = ? LIMIT 1`
		args = []any{hashURL(originalURL), originalURL}
	}
	row := r.db.QueryRowContext(ctx, query, args...)
	var u models.URL
	var createdAt dbTime
	if err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = createdAt.Time
	return &u, nil
}

func (r *sqlRepository) Create(ctx context.Context, originalURL, shortCode string) (*models.URL, error) {
	query := `INSERT INTO urls (original_url, short_code) VALUES (?, ?)`
	args := []any{originalURL, shortCode}
	if r.driver == "mysql" {
		query = `INSERT INTO urls (original_url, original_url_hash, short_code) VALUES (?, ?, ?)`
		args = []any{originalURL, hashURL(originalURL), shortCode}
	}
	res, err := r.db.ExecContext(ctx, query, args...)
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
	var createdAt dbTime
	if err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = createdAt.Time
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
		var createdAt dbTime
		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &createdAt); err != nil {
			return nil, err
		}
		u.CreatedAt = createdAt.Time
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

func hashURL(originalURL string) string {
	sum := sha256.Sum256([]byte(originalURL))
	return fmt.Sprintf("%x", sum)
}

type dbTime struct {
	time.Time
}

func (t *dbTime) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v.UTC()
		return nil
	case string:
		return t.parse(v)
	case []byte:
		return t.parse(string(v))
	default:
		return fmt.Errorf("unsupported database timestamp type %T", value)
	}
}

func (t *dbTime) parse(raw string) error {
	layouts := []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999", "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, raw, time.UTC); err == nil {
			t.Time = parsed.UTC()
			return nil
		}
	}
	return fmt.Errorf("invalid database timestamp %q", raw)
}
