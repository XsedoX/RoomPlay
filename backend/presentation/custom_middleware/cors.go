package customMiddleware

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
	"xsedox.com/main/config"
)

type CustomCors struct {
	corsMiddleware *cors.Cors
}

func NewCustomCors(configuration config.IConfiguration) *CustomCors {
	corsMiddleware := cors.New(cors.Options{
		AllowOriginVaryRequestFunc: func(r *http.Request, origin string) (bool, []string) {
			if strings.HasPrefix(r.RequestURI, "/api/v1/auth/google/callback") {
				return true, []string{}
			}
			return origin == configuration.Authentication().ClientOrigin, []string{}
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Device-Type", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: configuration.IsDevelopment(),
	})
	return &CustomCors{
		corsMiddleware: corsMiddleware,
	}
}

func (customCors *CustomCors) CorsHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request is for the Google callback and has no Origin, handle it and proceed.
			if strings.HasPrefix(r.URL.Path, "/api/v1/auth/google/callback") && r.Header.Get("Origin") == "" {
				next.ServeHTTP(w, r)
				return
			}
			// Otherwise, use the standard rs/cors handler.
			customCors.corsMiddleware.Handler(next).ServeHTTP(w, r)
		})
	}
}
