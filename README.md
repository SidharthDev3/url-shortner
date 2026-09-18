# URL Shortener

A modern full-stack URL shortener application built with a **Go** backend, a **React** (Vite) frontend, and **SQLite** persistence.

---

## Architecture Overview

```
url-shortner/
├── backend/          # Go RESTful API service & SQLite persistence
├── frontend/         # React (JavaScript) single-page application built with Vite
├── .gitignore        # Root gitignore for cross-environment hygiene
└── README.md         # Monorepo documentation & quickstart
```

- **[Backend](backend/README.md)**: Go modular architecture with clean separation across `cmd/`, `internal/`, `pkg/`, and `data/`.
- **[Frontend](frontend/README.md)**: React SPA scaffolded with Vite and organized into modular directories (`components/`, `pages/`, `services/`, `hooks/`, `context/`, `utils/`).
- **Database**: SQLite embedded database file stored under `backend/data/`.

---

## Quick Start (Prerequisites)

- [Go](https://go.dev/) (v1.22+)
- [Node.js](https://nodejs.org/) (v18+) & [npm](https://www.npmjs.com/)

### 1. Backend Setup
Navigate to the `backend/` directory:
```bash
cd backend
cp .env.example .env
go run cmd/server/main.go
```
See [backend/README.md](backend/README.md) for detailed configuration options.

### 2. Frontend Setup
Navigate to the `frontend/` directory:
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```
See [frontend/README.md](frontend/README.md) for detailed frontend documentation.

