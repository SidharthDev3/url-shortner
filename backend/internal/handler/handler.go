package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repository"
	"url-shortener-backend/internal/service"
)

// Handler holds dependencies for all HTTP endpoints.
type Handler struct {
	svc service.Service
}

// NewHandler creates a new HTTP Handler instance.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// ShortenHandler handles POST /api/shorten
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed: use POST")
		return
	}

	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	resp, err := h.svc.Shorten(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, service.ErrEmptyURL) ||
			errors.Is(err, service.ErrInvalidURL) ||
			errors.Is(err, service.ErrInvalidScheme) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// RedirectHandler handles GET /{shortCode}
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed: use GET")
		return
	}

	// Extract shortCode from path (e.g. "/abc123" -> "abc123")
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		writeJSONError(w, http.StatusNotFound, "short code is required")
		return
	}

	// Avoid colliding with API routes
	if strings.HasPrefix(path, "api/") {
		http.NotFound(w, r)
		return
	}

	originalURL, err := h.svc.Resolve(r.Context(), path)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "unknown short code: URL not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "database error")
		return
	}

	// Issue HTTP 302 Found redirection to original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// HistoryHandler handles GET /api/urls for the dashboard
func (h *Handler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed: use GET")
		return
	}

	history, err := h.svc.GetHistory(r.Context(), 20)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch URL history")
		return
	}

	if history == nil {
		history = []*models.ShortenResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(history)
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}
