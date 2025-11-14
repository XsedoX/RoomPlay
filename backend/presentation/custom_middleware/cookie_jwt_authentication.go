package customMiddleware

import (
	"context"
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/presentation/helpers"
	"xsedox.com/main/presentation/response"
)

type CookieJwtAuthentication struct {
	configuration config.IConfiguration
	jwtProvider   contracts.IJwtProvider
}

func NewCookieJwtAuthentication(configuration config.IConfiguration,
	jwtProvider contracts.IJwtProvider) *CookieJwtAuthentication {
	return &CookieJwtAuthentication{
		configuration: configuration,
		jwtProvider:   jwtProvider,
	}
}

func (jwtAuth *CookieJwtAuthentication) Next(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(helpers.RoomplayAccessTokenCookieName)
		if err != nil {
			response.WriteJsonFailure(w,
				"JwtAuthentication.MissingAuthCookie",
				"Missing cookie",
				err.Error(),
				r.URL.RequestURI(),
				http.StatusUnauthorized)
			return
		}

		userId, err := jwtAuth.jwtProvider.ValidateTokenAndGetUserId(authCookie.Value)
		if err != nil {
			response.WriteJsonFailure(w, "JwtAuthentication.TokenNotValid",
				"JWT issue",
				err.Error(),
				r.URL.RequestURI(),
				http.StatusUnauthorized)
			return
		}

		ctxWithClaims := context.WithValue(r.Context(), user.IdClaimContextKeyName, userId)
		next.ServeHTTP(w, r.WithContext(ctxWithClaims))
	})
}
