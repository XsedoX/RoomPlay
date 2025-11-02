package controllers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/infrastructure/authentication"
	"xsedox.com/main/presentation/response"
)

const (
	roomplayStateCookieExpirationTime    = time.Minute * 5
	stateCookieName                      = "roomPlay-state"
	RoomplayAccessTokenCookieName        = "roomplay-session-at"
	RoomplayRefreshTokenCookieName       = "roomplay-session-rt"
	roomPlayDeviceIdCookieName           = "roomplay-device-id"
	roomPlayDeviceIdCookieExpirationTime = 24 * time.Hour * 365 // a year
)

//TODO save where the user started logging and return to the same url

type OidcController struct {
	configuration             config.IConfiguration
	oidcAuthenticationService contracts.IOidcAuthenticationService
	googleOidcService         contracts.IGoogleOidcService
}

func NewOidcController(configuration config.IConfiguration,
	oidcAuthenticationService contracts.IOidcAuthenticationService,
	googleOidcService contracts.IGoogleOidcService) *OidcController {
	return &OidcController{
		configuration:             configuration,
		oidcAuthenticationService: oidcAuthenticationService,
		googleOidcService:         googleOidcService,
	}
}

// HandleLoginWithGoogle godoc
//
//	@Summary		Initiate Google OAuth2 login
//	@Description	Starts the OAuth2 flow by generating a state cookie and redirecting to Google's authorization endpoint.
//	@Tags			AuthenticationController
//	@Accept			json
//	@Produce		json
//	@Success		302	{string}	string	"Redirect to Google OAuth"
//	@Failure		403	{object}	response.Error	"Invalid redirect URL"
//	@Router			/auth/google/login [post]
func (handler *OidcController) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	clearStateCookie(w)
	state := setStateCookie(w, handler.configuration.Server().BasePath)

	googleUrl, err := handler.googleOidcService.GenerateOidcUrl(state)
	if err != nil {
		response.WriteJsonFailure(w, "Couldn't generate google oauth url", http.StatusForbidden)
	}

	response.WriteJsonSuccess(w, googleUrl, http.StatusOK)
}
func (handler *OidcController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		response.WriteJsonFailure(w, "Problem with parsing query", http.StatusInternalServerError)
		return
	}
	if !verifyStateCookie(r, params.Get("state")) {
		response.WriteJsonFailure(w, "Invalid state", http.StatusForbidden)
		return
	}

	code := params.Get("code")
	if code == "" {
		response.WriteJsonFailure(w, "Invalid code", http.StatusForbidden)
		return
	}

	var deviceIdValue *shared.DeviceId
	deviceId, err := r.Cookie(roomPlayDeviceIdCookieName)
	if err != nil {
		deviceIdValue = nil
	} else {
		deviceIdValue = shared.ParseDeviceId(deviceId.Value)
	}

	apiTokenResponse, err := handler.oidcAuthenticationService.AuthenticateWithGoogle(r.Context(), code, deviceIdValue)
	if err != nil {
		response.WriteJsonFailure(w, err.Error(), http.StatusForbidden)
		return
	}

	setAccessTokenCookie(w, apiTokenResponse.AccessToken, handler.configuration.Server().BasePath)
	setRefreshTokenCookie(w, apiTokenResponse.RefreshToken, handler.configuration.Server().BasePath)
	setDeviceIdCookie(w, apiTokenResponse.DeviceId.String(), handler.configuration.Server().BasePath)

	http.Redirect(w, r, handler.configuration.Authentication().ClientOrigin+"/signin-oidc", http.StatusPermanentRedirect)
}
func setAccessTokenCookie(w http.ResponseWriter, accessToken, basePath string) {
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
func setRefreshTokenCookie(w http.ResponseWriter, refreshToken, basePath string) {
	expiresAt := time.Now().Add(credentials.RefreshTokenExpirationTime).UTC()
	http.SetCookie(w, &http.Cookie{
		Name:     RoomplayRefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  expiresAt,
		Path:     basePath + "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setDeviceIdCookie(w http.ResponseWriter, deviceId string, basePath string) {
	expiresAt := time.Now().UTC().Add(roomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     roomPlayDeviceIdCookieName,
		Value:    deviceId,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setStateCookie(w http.ResponseWriter, basePath string) string {
	expiresAt := time.Now().UTC().Add(roomplayStateCookieExpirationTime)
	state := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Expires:  expiresAt,
		MaxAge:   int(roomplayStateCookieExpirationTime.Seconds()),
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}
func clearStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   stateCookieName,
		MaxAge: -1,
	})
}
func verifyStateCookie(r *http.Request, stateFromUrl string) bool {
	state, err := r.Cookie(stateCookieName)
	if err != nil {
		return false
	}
	stateString := state.Value

	if stateFromUrl != stateString {
		return false
	}
	return true
}
