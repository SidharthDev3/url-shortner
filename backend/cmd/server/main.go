package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/handler"
	"url-shortener-backend/internal/middleware"
	"url-shortener-backend/internal/repository"
	"url-shortener-backend/internal/service"
)

func main() {
	// 1. Load application configuration
	cfg := config.Load()
	log.Printf("Starting URL Shortener Backend (Env: %s, Port: %s)...", cfg.Env, cfg.Port)

	// 2. Initialize SQLite database repository
	repo, err := repository.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer repo.Close()
	log.Printf("SQLite database connected successfully at: %s", cfg.DBPath)

	// 3. Initialize business service & HTTP handler
	svc := service.NewURLService(repo, cfg.BaseURL)
	h := handler.NewHandler(svc)

	// 4. Set up HTTP routing
	mux := http.NewServeMux()
	mux.HandleFunc("/api/shorten", h.ShortenHandler)
	mux.HandleFunc("/api/urls", h.HistoryHandler)
	mux.HandleFunc("/", h.RedirectHandler)

	// 5. Wrap router with CORS middleware
	handlerWithCORS := middleware.EnableCORS(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handlerWithCORS,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Start server in goroutine
	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	// 7. Wait for graceful shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown forced: %v", err)
	}

	log.Println("Server exiting.")
}
