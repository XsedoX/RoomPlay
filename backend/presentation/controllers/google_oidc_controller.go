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
	"xsedox.com/main/presentation/presentationErrors"
)

const (
	roomplayStateCookieExpirationTime = time.Minute * 5
	stateCookieName                   = "roomPlay-state"
	roomplayAccessTokenCookieName     = "roomplay-session-at"
	roomplayRefreshTokenCookieName    = "roomplay-session-rt"
	roomPlayDeviceIdCookieName        = "roomplay-device-id"
	roomPlayDeviceIdCookieLifespan    = 24 * time.Hour * 365 // a year
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
//	@Router			/api/v1/auth/google/login [post]
func (handler *OidcController) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	clearStateCookie(w)
	state := setStateCookie(w)

	googleUrl, err := handler.googleOidcService.GenerateOidcUrl(state)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Couldnt generate google oauth url", http.StatusForbidden)
	}

	presentationErrors.WriteJsonSuccess(w, googleUrl, http.StatusOK)
}
func (handler *OidcController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Problem with parsing query", http.StatusInternalServerError)
		return
	}
	if !verifyStateCookie(r, params.Get("state")) {
		presentationErrors.WriteJsonFailure(w, "Invalid state", http.StatusForbidden)
		return
	}

	code := params.Get("code")
	if code == "" {
		presentationErrors.WriteJsonFailure(w, "Invalid code", http.StatusForbidden)
		return
	}

	var deviceIdValue *shared.DeviceId
	deviceId, err := r.Cookie(roomPlayDeviceIdCookieName)
	if err != nil {
		deviceIdValue = nil
	} else {
		deviceIdValue, err = shared.ParseDeviceId(deviceId.Value)
		if err != nil {
			deviceIdValue = nil
		}
	}

	apiTokenResponse, err := handler.oidcAuthenticationService.AuthenticateWithGoogle(r.Context(), code, deviceIdValue)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "couldn't authenticate", http.StatusForbidden)
		return
	}

	setAccessTokenCookie(w, apiTokenResponse.AccessToken)
	setRefreshTokenCookie(w, apiTokenResponse.RefreshToken)
	setDeviceIdCookie(w, apiTokenResponse.DeviceId.String())

	http.Redirect(w, r, handler.configuration.Authentication().ClientOrigin+"/signin-oidc", http.StatusPermanentRedirect)
}
func setAccessTokenCookie(w http.ResponseWriter, accessToken string) {
	expiresAt := time.Now().UTC().Add(authentication.AccessTokenExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     roomplayAccessTokenCookieName,
		Value:    accessToken,
		Expires:  expiresAt,
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	expiresAt := time.Now().Add(credentials.RefreshTokenExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     roomplayRefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  expiresAt,
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setDeviceIdCookie(w http.ResponseWriter, deviceId string) {
	expiresAt := time.Now().UTC().Add(roomPlayDeviceIdCookieLifespan)
	http.SetCookie(w, &http.Cookie{
		Name:     roomPlayDeviceIdCookieName,
		Value:    deviceId,
		Expires:  expiresAt,
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setStateCookie(w http.ResponseWriter) string {
	expiresAt := time.Now().UTC().Add(roomplayStateCookieExpirationTime)
	state := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Expires:  expiresAt,
		MaxAge:   int(roomplayStateCookieExpirationTime.Seconds()),
		Path:     "/api/v1",
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
