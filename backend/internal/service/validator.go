package service

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrEmptyURL   = errors.New("url cannot be empty")
	ErrInvalidURL = errors.New("invalid url format: must include a valid domain")
	ErrInvalidScheme = errors.New("invalid protocol: url must start with http:// or https://")
)

// ValidateURL validates that the raw input is a legitimate HTTP/HTTPS web address.
func ValidateURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", ErrEmptyURL
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return "", ErrInvalidURL
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidScheme
	}

	if parsed.Host == "" {
		return "", ErrInvalidURL
	}

	return parsed.String(), nil
}

