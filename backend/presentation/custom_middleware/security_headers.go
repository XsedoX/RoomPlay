package customMiddleware

import (
	"net/http"

	"xsedox.com/main/config"
)

type SecurityHeaders struct {
	configuration config.IConfiguration
}

func NewSecurityHeaders(configuration config.IConfiguration) *SecurityHeaders {
	return &SecurityHeaders{configuration: configuration}
}

func (s *SecurityHeaders) Next(next http.Handler) http.Handler {
	cspValue := "default-src 'none';frame-ancestors 'none';"
	isReportOnly := s.configuration.IsDevelopment()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isReportOnly {
			w.Header().Set("Content-Security-Policy-Report-Only", cspValue)
		} else {
			w.Header().Set("Content-Security-Policy", cspValue)
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")

		next.ServeHTTP(w, r)
	})
}
