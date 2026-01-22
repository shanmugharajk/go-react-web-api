# go-react-web-api

Production-style Go (Chi, GORM) and React monorepo demonstrating clean architecture, API design, and real-world patterns.

## Project Structure

```
.
├── api/          # Go backend (REST API)
└── web/          # React frontend
```

## Backend (api/)

Clean, idiomatic Go backend with domain-oriented architecture.

### Tech Stack

- **Go 1.25.5**
- **Router**: chi v5.2.4 - Lightweight, idiomatic HTTP router
- **ORM**: GORM v1.31.1 - Developer-friendly ORM
- **Database**: SQLite 3 (local development)
- **Logging**: log/slog (standard library)

### Quick Start

```bash
cd api
go run cmd/server/main.go
```

Server starts at `http://localhost:8080`

See `api/README.md` for detailed documentation.
