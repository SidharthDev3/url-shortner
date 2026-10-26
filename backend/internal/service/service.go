package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repository"
	"url-shortener-backend/pkg/utils"
)

var (
	ErrCollision = errors.New("failed to generate unique short code, please try again")
)

// Service defines high-level business operations for URL shortening.
type Service interface {
	Shorten(ctx context.Context, rawURL string) (*models.ShortenResponse, error)
	Resolve(ctx context.Context, shortCode string) (string, error)
	GetHistory(ctx context.Context, limit int) ([]*models.ShortenResponse, error)
}

type URLService struct {
	repo    repository.Repository
	baseURL string
}

// NewURLService creates a new URLService instance.
func NewURLService(repo repository.Repository, baseURL string) *URLService {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &URLService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// Shorten processes a long URL: validates, checks for duplicates, creates a short code, and persists it.
func (s *URLService) Shorten(ctx context.Context, rawURL string) (*models.ShortenResponse, error) {
	// 1. Validate URL (Report Section 3.4 & 4.2.2)
	validURL, err := ValidateURL(rawURL)
	if err != nil {
		return nil, err
	}

	// 2. Duplicate Check (Report Section 5.3 & 5.4: if mapping exists, return existing short URL)
	existing, err := s.repo.FindByOriginalURL(ctx, validURL)
	if err == nil && existing != nil {
		return &models.ShortenResponse{
			ShortCode:   existing.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, existing.ShortCode),
			OriginalURL: existing.OriginalURL,
			CreatedAt:   existing.CreatedAt,
		}, nil
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("error querying database for duplicates: %w", err)
	}

	// 3. Generate unique short code with collision handling (crypto/rand)
	var shortCode string
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		candidate, err := utils.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Ensure candidate code does not already exist
		_, err = s.repo.FindByShortCode(ctx, candidate)
		if errors.Is(err, repository.ErrNotFound) {
			shortCode = candidate
			break
		}
	}

	if shortCode == "" {
		return nil, ErrCollision
	}

	// 4. Save to SQLite database
	saved, err := s.repo.Create(ctx, validURL, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to save url to database: %w", err)
	}

	return &models.ShortenResponse{
		ShortCode:   saved.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, saved.ShortCode),
		OriginalURL: saved.OriginalURL,
		CreatedAt:   saved.CreatedAt,
	}, nil
}

// Resolve retrieves the original long URL given a short code.
func (s *URLService) Resolve(ctx context.Context, shortCode string) (string, error) {
	shortCode = strings.TrimSpace(shortCode)
	if shortCode == "" {
		return "", repository.ErrNotFound
	}

	u, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return u.OriginalURL, nil
}

// GetHistory returns recently shortened links for the dashboard view.
func (s *URLService) GetHistory(ctx context.Context, limit int) ([]*models.ShortenResponse, error) {
	urls, err := s.repo.GetRecentURLs(ctx, limit)
	if err != nil {
		return nil, err
	}

	res := make([]*models.ShortenResponse, 0, len(urls))
	for _, u := range urls {
		res = append(res, &models.ShortenResponse{
			ShortCode:   u.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, u.ShortCode),
			OriginalURL: u.OriginalURL,
			CreatedAt:   u.CreatedAt,
		})
	}

	return res, nil
}
