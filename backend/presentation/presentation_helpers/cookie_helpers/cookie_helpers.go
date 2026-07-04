package cookie_helpers

import (
	"net/http"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/jwt_provider"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
)

func ClearAccessTokenCookie(w http.ResponseWriter, basePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayAccessTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearRefreshTokenCookie(w http.ResponseWriter, basePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     basePath + constants.RefreshTokenPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAccessTokenCookie(w http.ResponseWriter, accessToken, basePath string) {
	expiresAt := time.Now().Add(jwt_provider.AccessTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayAccessTokenCookieName,
		Value:    accessToken,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken, basePath string) {
	expiresAt := time.Now().Add(internal_credentials.RefreshTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  expiresAt,
		Path:     basePath + constants.RefreshTokenPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetDeviceIdCookie(w http.ResponseWriter, deviceId string, basePath string) {
	expiresAt := time.Now().UTC().Add(constants.RoomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayDeviceIdCookieName,
		Value:    deviceId,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
