# Backend - URL Shortener

Go RESTful API service for the URL Shortener application with SQLite persistence.

---

## Project Layout

This backend adheres to standard Go project conventions:

```
backend/
├── .env.example          # Template for environment configuration
├── .gitignore            # Go-specific gitignore rules
├── README.md             # Backend architecture and instructions
├── go.mod                # Go module definition
├── cmd/
│   └── server/
│       └── main.go       # Application entry point (server bootstrap)
├── data/
│   └── .gitkeep          # SQLite database storage directory
├── internal/             # Private application packages (not importable by other apps)
│   ├── config/           # Environment and configuration loaders
│   ├── handler/          # HTTP request handlers & routing
│   ├── middleware/       # HTTP middleware (logging, recovery, CORS, rate limiting)
│   ├── models/           # Domain entities & data transfer objects (DTOs)
│   ├── repository/       # Data access layer & SQLite database interactions
│   └── service/          # Core business logic
└── pkg/                  # Public shared packages & utilities
    └── utils/            # General-purpose helpers (e.g. base62 encoding)
```

---

## Folder Purposes

- **`cmd/server/`**: Contains the `main.go` entry point. It handles flags, config initialization, dependency injection, and starts the HTTP server.
- **`internal/config/`**: Loads configuration from environment variables or `.env` files into strongly-typed structs.
- **`internal/handler/`**: Implements HTTP endpoints, parses incoming JSON/requests, and writes responses.
- **`internal/middleware/`**: Handles cross-cutting concerns such as CORS headers, request logging, authentication, and error recovery.
- **`internal/models/`**: Defines core domain models (e.g., `URL`, `Analytics`) and request/response payloads.
- **`internal/repository/`**: Abstraction and implementation for SQLite data persistence and SQL queries.
- **`internal/service/`**: Encapsulates business rules (e.g. short code generation, URL validation, expiration checks).
- **`pkg/utils/`**: Reusable helper functions that have no dependencies on internal domain logic.
- **`data/`**: Designated runtime storage directory for SQLite database files (`*.db`).

---

## Getting Started

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

