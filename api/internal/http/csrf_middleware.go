package http

import (
	"net/http"
	"strings"
)

// csrfProtectUnlessBearerToken wraps a CSRF protection middleware and exempts requests
// that use Bearer token authentication.
//
// Security rationale:
// - CSRF attacks exploit the browser's automatic cookie sending behavior
// - Bearer tokens in Authorization headers are NOT automatically sent by browsers
// - Therefore, requests with Bearer tokens are not vulnerable to CSRF
// - This allows API clients (Postman, mobile apps, etc.) to bypass CSRF complexity
//
// IMPORTANT: This ONLY exempts Bearer token auth. Cookie-based auth MUST still use CSRF.
func csrfProtectUnlessBearerToken(csrfMiddleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if request has Bearer token in Authorization header
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				// Skip CSRF protection for Bearer token requests
				// These are not vulnerable to CSRF attacks
				next.ServeHTTP(w, r)
				return
			}

			// Apply CSRF protection for all other requests
			// (including cookie-based auth and unauthenticated requests)
			csrfMiddleware(next).ServeHTTP(w, r)
		})
	}
}
