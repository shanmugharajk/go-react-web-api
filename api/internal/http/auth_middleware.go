package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/jwt"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/sessions"
)

// RequireAuth is a middleware that validates authentication (session or JWT) and injects user and user ID into context.
// It supports dual authentication: cookie-based sessions OR JWT bearer tokens.
func RequireAuth(sessionStore *sessions.SQLiteStore, jwtService *jwt.TokenService, authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID int64
			var authenticated bool

			// Try JWT Bearer token first (for API clients)
			if authHeader := r.Header.Get("Authorization"); authHeader != "" {
				// Check for Bearer token
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenString := strings.TrimPrefix(authHeader, "Bearer ")

					// Validate JWT token
					id, err := jwtService.Validate(tokenString)
					if err == nil {
						userID = id
						authenticated = true
					}
					// If JWT validation fails, don't fall through to session auth
					// This prevents mixing auth methods in a single request
					if !authenticated {
						if errors.Is(err, jwt.ErrExpiredToken) {
							response.Error(w, http.StatusUnauthorized, "token expired")
						} else {
							response.Error(w, http.StatusUnauthorized, "invalid token")
						}
						return
					}
				}
			}

			// Try session cookie (for browser/SPA clients) if JWT didn't authenticate
			if !authenticated {
				cookie, err := r.Cookie("session")
				if err != nil {
					response.Error(w, http.StatusUnauthorized, "unauthorized")
					return
				}

				// Validate session
				session, err := sessionStore.Get(cookie.Value)
				if err != nil {
					if errors.Is(err, sessions.ErrSessionExpired) {
						response.Error(w, http.StatusUnauthorized, "session expired")
					} else {
						response.Error(w, http.StatusUnauthorized, "invalid session")
					}
					return
				}

				// Check if session has expired (redundant but explicit)
				if session.IsExpired() {
					response.Error(w, http.StatusUnauthorized, "session expired")
					return
				}

				userID = session.UserID
				authenticated = true
			}

			// At this point, we must be authenticated
			if !authenticated {
				response.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// Fetch the full user object from the database
			user, err := authService.GetUserByID(uint(userID))
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
				return
			}

			// Inject both user ID and user object into context (type-safe)
			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.UserIDKey, userID)
			ctx = context.WithValue(ctx, auth.UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
