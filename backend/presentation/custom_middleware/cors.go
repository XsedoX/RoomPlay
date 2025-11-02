package customMiddleware

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

func CorsHandler(c *cors.Cors) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request is for the Google callback and has no Origin, handle it and proceed.
			if strings.HasPrefix(r.URL.Path, "/api/v1/auth/google/callback") && r.Header.Get("Origin") == "" {
				next.ServeHTTP(w, r)
				return
			}
			// Otherwise, use the standard rs/cors handler.
			c.Handler(next).ServeHTTP(w, r)
		})
	}
}
