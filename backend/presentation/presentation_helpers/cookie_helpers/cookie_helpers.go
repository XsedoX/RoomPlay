package cookie_helpers

import (
	"net/http"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/jwt_provider"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/google/uuid"
)

const accessTokenCookiePath = constants.ApiBasePath

func ClearAccessTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayAccessTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     accessTokenCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAccessTokenCookie(w http.ResponseWriter, accessToken string) {
	expiresAt := time.Now().Add(jwt_provider.AccessTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayAccessTokenCookieName,
		Value:    accessToken,
		Expires:  expiresAt,
		Path:     accessTokenCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     constants.RefreshTokenCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	expiresAt := time.Now().Add(internal_credentials.RefreshTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  expiresAt,
		Path:     constants.RefreshTokenCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetDeviceIdCookie(w http.ResponseWriter, deviceId string) {
	expiresAt := time.Now().UTC().Add(constants.RoomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayDeviceIdCookieName,
		Value:    deviceId,
		Expires:  expiresAt,
		Path:     constants.ApiBasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetDeviceTypeCookie(w http.ResponseWriter, deviceType string) {
	expiresAt := time.Now().UTC().Add(constants.RoomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayDeviceTypeCookieName,
		Value:    deviceType,
		Expires:  expiresAt,
		Path:     constants.ApiBasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetStateCookie(w http.ResponseWriter) string {
	expiresAt := time.Now().Add(constants.RoomPlayStateCookieExpirationTime).UTC()
	state := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayStateCookieName,
		Value:    state,
		Expires:  expiresAt,
		MaxAge:   int(constants.RoomPlayStateCookieExpirationTime.Seconds()),
		Path:     constants.ApiBasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}

func ClearStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RoomPlayStateCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     constants.ApiBasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func VerifyStateCookie(r *http.Request, stateFromUrl string) bool {
	state, err := r.Cookie(constants.RoomPlayStateCookieName)
	if err != nil {
		return false
	}
	stateString := state.Value

	return stateString == stateFromUrl
}
