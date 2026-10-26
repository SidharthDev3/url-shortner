package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"url-shortener-backend/internal/repository"
)

func setupTestRepo(t *testing.T) (*repository.SQLiteRepository, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "url_shortener_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to init test repo: %v", err)
	}

	cleanup := func() {
		repo.Close()
		os.RemoveAll(tempDir)
	}

	return repo, cleanup
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid HTTPS", "https://example.com/test", false},
		{"Valid HTTP", "http://localhost:3000/path?query=1", false},
		{"Invalid scheme ftp", "ftp://example.com", true},
		{"Missing scheme", "example.com", true},
		{"Empty string", "", true},
		{"Spaces only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) err = %v, wantErr = %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestURLService_ShortenAndDuplicateCheck(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	svc := NewURLService(repo, "http://localhost:8080")
	ctx := context.Background()

	targetURL := "https://example.com/long-page-for-testing"

	// 1. Initial shorten
	resp1, err := svc.Shorten(ctx, targetURL)
	if err != nil {
		t.Fatalf("unexpected error on first shorten: %v", err)
	}

	if resp1.ShortCode == "" {
		t.Errorf("expected non-empty short code")
	}
	if len(resp1.ShortCode) != 6 {
		t.Errorf("expected short code of length 6, got %d", len(resp1.ShortCode))
	}
	if resp1.OriginalURL != targetURL {
		t.Errorf("expected %s, got %s", targetURL, resp1.OriginalURL)
	}

	// 2. Duplicate Check: Shortening the same URL must return the identical short code
	resp2, err := svc.Shorten(ctx, targetURL)
	if err != nil {
		t.Fatalf("unexpected error on duplicate shorten: %v", err)
	}

	if resp2.ShortCode != resp1.ShortCode {
		t.Errorf("duplicate check failed: expected short code %s, got %s", resp1.ShortCode, resp2.ShortCode)
	}

	// 3. Resolve short code
	resolved, err := svc.Resolve(ctx, resp1.ShortCode)
	if err != nil {
		t.Fatalf("failed to resolve short code: %v", err)
	}
	if resolved != targetURL {
		t.Errorf("resolved URL mismatch: expected %s, got %s", targetURL, resolved)
	}

	// 4. Resolve unknown code
	_, err = svc.Resolve(ctx, "unknownXYZ")
	if err == nil {
		t.Errorf("expected error when resolving non-existent short code")
	}
}

