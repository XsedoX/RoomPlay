package customMiddleware

import (
	"context"
	"encoding/base64"
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
				"The authentication cookie has not been sent",
				r.URL.RequestURI(),
				http.StatusUnauthorized)
			return
		}
		decodedToken, err := base64.StdEncoding.DecodeString(authCookie.Value)
		if err != nil {
			// Handle malformed cookie value.
			response.WriteJsonFailure(w,
				"JwtAuthentication.DecodeString",
				"Invalid access token",
				"The JWT token could not be decoded",
				r.URL.RequestURI(),
				http.StatusUnauthorized)
			return
		}

		userId, err := jwtAuth.jwtProvider.ValidateTokenAndGetUserId(string(decodedToken))
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
