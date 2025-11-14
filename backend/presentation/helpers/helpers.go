package helpers

import (
	"net/http"
	"time"

	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/infrastructure/authentication"
)

func ClearAccessTokenCookie(w http.ResponseWriter, basePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     RoomplayAccessTokenCookieName,
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
		Name:     RoomplayRefreshTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     basePath + "/auth/refresh-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAccessTokenCookie(w http.ResponseWriter, accessToken, basePath string) {
	expiresAt := time.Now().Add(authentication.AccessTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     RoomplayAccessTokenCookieName,
		Value:    accessToken,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken, basePath string) {
	expiresAt := time.Now().Add(credentials.RefreshTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     RoomplayRefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  expiresAt,
		Path:     basePath + "/auth/refresh-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
