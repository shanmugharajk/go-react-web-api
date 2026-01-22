# Go React Web API

A full-stack web application with Go backend API and React frontend.


## Project Structure

```
.
├── api/          # Go backend (REST API)
└── web/          # React frontend
```

## Tech Stack

### Backend
- **Go** 1.25.5
- **Chi** v5.2.4 - HTTP router
- **GORM** v1.31.1 - ORM
- **SQLite** - Database
- **golang.org/x/crypto** - Password hashing (Argon2id)
- **gorilla/csrf** v1.7.3 - CSRF protection

### Frontend
- Coming soon

## Getting Started

```bash
cd api

# Install dependencies
go mod download

# Run the server
go run ./cmd/server/main.go

# Server starts on http://localhost:8080
```

Check by visiting health check endpoint - `GET /healthz`.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `DATABASE_DSN` | `file:./data/pos.db` | SQLite database path |
| `AUTH_SESSION_SECRET` | - | HMAC secret for sessions (generate with `openssl rand -hex 32`) |
| `AUTH_CSRF_SECRET` | - | CSRF token secret (generate with `openssl rand -hex 32`) |
| `IS_DEVELOPMENT` | `true` | Development mode (set to `false` in production) |
| `AUTH_SESSION_DURATION` | `86400` | Session duration in seconds (default: 24 hours) |

## License

MIT
