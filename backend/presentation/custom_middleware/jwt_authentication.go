package customMiddleware

import (
	"context"
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/presentation/presentationErrors"
)

type JwtAuthentication struct {
	configuration config.IConfiguration
	jwtProvider   contracts.IJwtProvider
}

func NewJwtAuthentication(configuration config.IConfiguration,
	jwtProvider contracts.IJwtProvider) *JwtAuthentication {
	return &JwtAuthentication{
		configuration: configuration,
		jwtProvider:   jwtProvider,
	}
}

func (jwtAuth *JwtAuthentication) Next(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			presentationErrors.WriteJsonFailure(w, "Missing Auth", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		userId, err := jwtAuth.jwtProvider.ValidateTokenAndGetUserId(tokenString)
		if err != nil {
			presentationErrors.WriteJsonFailure(w, "JWT issue", http.StatusUnauthorized)
			return
		}

		ctxWithClaims := context.WithValue(r.Context(), user.IdClaimContextKeyName, userId)
		next.ServeHTTP(w, r.WithContext(ctxWithClaims))
	})
}
