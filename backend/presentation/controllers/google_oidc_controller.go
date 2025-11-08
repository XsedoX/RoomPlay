package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/presentation/helpers"
	"xsedox.com/main/presentation/response"
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
	state := setStateCookie(w, handler.configuration.Server().BasePath)
	deviceType := r.Header.Get("X-Device-Type")
	parsedDeviceType := user.ParseDeviceType(&deviceType)
	if parsedDeviceType == nil {
		response.WriteJsonFailure(w,
			"OidcController.ParseDeviceType",
			"Device type is not valid",
			fmt.Sprintf("Device type should be one of the following values: %s", strings.Join(user.ListDeviceTypes(), ", ")),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}
	setDeviceTypeCookie(w, parsedDeviceType.String(), handler.configuration.Server().BasePath)

	googleUrl, err := handler.googleOidcService.GenerateOidcUrl(state)
	if err != nil {
		response.WriteJsonFailure(w,
			"OidcController.GenerateOidcUrl",
			"Something went wrong with generating google url",
			err.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}

	response.WriteJsonSuccess(w, googleUrl, http.StatusOK)
}
func (handler *OidcController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		response.WriteJsonFailure(w,
			"OidcController.ParseQuery",
			"Problem with parsing query",
			err.Error(),
			r.URL.RequestURI(),
			http.StatusInternalServerError)
		return
	}
	if !verifyStateCookie(r, params.Get("state")) {
		response.WriteJsonFailure(w,
			"OidcController.VerifyStateCookie",
			"Invalid state",
			"State cookie is not valid",
			r.URL.RequestURI(),
			http.StatusForbidden)
		return
	}

	code := params.Get("code")
	if code == "" {
		response.WriteJsonFailure(w,
			"OidcController.GetCode",
			"Invalid code",
			"Code is empty",
			r.URL.RequestURI(),
			http.StatusForbidden)
		return
	}

	deviceType, err := r.Cookie(helpers.RoomPlayDeviceTypeCookieName)
	if err != nil {
		response.WriteJsonFailure(w,
			"OidcController.GetDeviceTypeCookie",
			"Couldn't get device type cookie",
			err.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}

	var deviceIdValue *user.DeviceId
	deviceId, err := r.Cookie(helpers.RoomPlayDeviceIdCookieName)
	if err != nil {
		deviceIdValue = nil
	} else {
		deviceIdValue = user.ParseDeviceId(deviceId.Value)
	}

	apiTokenResponse, err := handler.oidcAuthenticationService.AuthenticateWithGoogle(r.Context(), code, deviceIdValue, user.ParseDeviceType(&deviceType.Value))
	if err != nil {
		response.WriteJsonApplicationFailure(w, err, r.URL.RequestURI())
		return
	}
	base64AccessToken := base64.StdEncoding.EncodeToString([]byte(apiTokenResponse.AccessToken))
	base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(apiTokenResponse.RefreshToken))
	helpers.SetAccessTokenCookie(w, base64AccessToken, handler.configuration.Server().BasePath)
	helpers.SetRefreshTokenCookie(w, base64RefreshToken, handler.configuration.Server().BasePath)
	setDeviceIdCookie(w, *apiTokenResponse.DeviceId.String(), handler.configuration.Server().BasePath)
	clearStateCookie(w, handler.configuration.Server().BasePath)

	http.Redirect(w, r, handler.configuration.Authentication().ClientOrigin+"/signin-oidc", http.StatusPermanentRedirect)
}
func setDeviceTypeCookie(w http.ResponseWriter, deviceType string, basePath string) {
	expiresAt := time.Now().UTC().Add(helpers.RoomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     helpers.RoomPlayDeviceTypeCookieName,
		Value:    deviceType,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setDeviceIdCookie(w http.ResponseWriter, deviceId string, basePath string) {
	expiresAt := time.Now().UTC().Add(helpers.RoomPlayDeviceIdCookieExpirationTime)
	http.SetCookie(w, &http.Cookie{
		Name:     helpers.RoomPlayDeviceIdCookieName,
		Value:    deviceId,
		Expires:  expiresAt,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func setStateCookie(w http.ResponseWriter, basePath string) string {
	expiresAt := time.Now().Add(helpers.RoomplayStateCookieExpirationTime).UTC()
	state := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     helpers.RoomplayStateCookieName,
		Value:    state,
		Expires:  expiresAt,
		MaxAge:   int(helpers.RoomplayStateCookieExpirationTime.Seconds()),
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}
func clearStateCookie(w http.ResponseWriter, basePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     helpers.RoomplayStateCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     basePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
func verifyStateCookie(r *http.Request, stateFromUrl string) bool {
	state, err := r.Cookie(helpers.RoomplayStateCookieName)
	if err != nil {
		return false
	}
	stateString := state.Value

	if stateFromUrl != stateString {
		return false
	}
	return true
}
