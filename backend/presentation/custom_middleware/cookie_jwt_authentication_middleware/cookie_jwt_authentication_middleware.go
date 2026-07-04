package cookie_jwt_authentication_middleware

import (
	"context"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

type CookieJwtAuthentication struct {
	configuration config.IConfiguration
	jwtProvider   i_jwt_provider.IJwtProvider
}

func NewCookieJwtAuthentication(configuration config.IConfiguration,
	jwtProvider i_jwt_provider.IJwtProvider,
) *CookieJwtAuthentication {
	return &CookieJwtAuthentication{
		configuration: configuration,
		jwtProvider:   jwtProvider,
	}
}

func (jwtAuth *CookieJwtAuthentication) Next(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(constants.RoomPlayAccessTokenCookieName)
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
