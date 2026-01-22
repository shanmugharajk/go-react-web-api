# Go React Web API - Backend

A clean, idiomatic Go backend with domain-oriented architecture.

## Structure

```
api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── app/            # Application bootstrap
│   ├── config/         # Configuration management
│   ├── db/             # Database setup
│   ├── http/           # HTTP server and routing
│   ├── modules/        # Domain modules (auth, product)
│   └── pkg/            # Shared utilities
```

## Getting Started

### Prerequisites

- Go 1.23 or later
- SQLite3

### Installation

1. Install dependencies:
```bash
cd api
go mod download
```

2. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.
`GET /healthz` - Health check endpoint

### Environment Variables

Configure via `.env` file or environment variables:

- `PORT` - HTTP server port (default: 8080)
- `HOST` - HTTP server host (default: localhost)
- `DATABASE_DSN` - Database connection string (default: file:./data/app.db)

## Development

### Project Principles

- Idiomatic Go code
- Clear separation of concerns
- Domain-oriented modules
- Minimal framework magic
- Easy to test and reason about

### Module Structure

Each domain module follows a consistent pattern:

- `handler.go` - HTTP layer
- `service.go` - Business logic
- `repository.go` - Data access
- `model.go` - Domain models

## Production

### Building

Build the application for production:

```bash
cd api
go build -o bin/server cmd/server/main.go
```

### Running in Production

1. Set production environment variables:
```bash
export PORT=8080
export HOST=0.0.0.0
export DATABASE_DSN="file:/path/to/prod/database.db"
```

2. Run the compiled binary:
```bash
./bin/server
```

### Production Considerations

- Use environment variables instead of `.env` files
- Configure a production-ready database (PostgreSQL, MySQL, etc.) - Used sqlite for simplicity now
- Set appropriate logging levels for production
- Use a process manager like systemd or Docker for deployment - Directly deployed in Digital ocean droplet for simplicity as of now
- Enable HTTPS in production with proper certificates

## Tech Stack

- **Router**: [chi](https://github.com/go-chi/chi) - Lightweight, idiomatic HTTP router
- **ORM**: [GORM](https://gorm.io) - Developer-friendly ORM
- **Database**: SQLite (for local development)
- **Logging**: `log/slog` (standard library)
