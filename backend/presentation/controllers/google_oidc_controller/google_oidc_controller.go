package google_oidc_controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_google_oidc_service"
	"github.com/XsedoX/RoomPlay/application/services/services_contracts/i_oidc_authentication_service"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/cookie_helpers"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

// TODO: save where the user started logging and return to the same url
type GoogleOidcController struct {
	configuration             config.IConfiguration
	oidcAuthenticationService i_oidc_authentication_service.IOidcAuthenticationService
	googleOidcService         i_google_oidc_service.IGoogleOidcService
}

func NewOidcController(configuration config.IConfiguration,
	oidcAuthenticationService i_oidc_authentication_service.IOidcAuthenticationService,
	googleOidcService i_google_oidc_service.IGoogleOidcService,
) *GoogleOidcController {
	return &GoogleOidcController{
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
func (handler *GoogleOidcController) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	state := cookie_helpers.SetStateCookie(w)
	deviceType := r.Header.Get("X-Device-Type")
	parsedDeviceType := device_type.ParseDeviceType(&deviceType)
	if parsedDeviceType == nil {
		response.WriteJsonFailure(w,
			"OidcController.ParseDeviceType",
			"Device type is not valid",
			fmt.Sprintf("Device type should be one of the following values: %s", strings.Join(device_type.ListDeviceTypes(), ", ")),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}
	cookie_helpers.SetDeviceTypeCookie(w, parsedDeviceType.String())

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

	response.WriteJsonSuccess(w, googleUrl)
}

func (handler *GoogleOidcController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
	if !cookie_helpers.VerifyStateCookie(r, params.Get("state")) {
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

	deviceType, err := r.Cookie(constants.RoomPlayDeviceTypeCookieName)
	if err != nil {
		response.WriteJsonFailure(w,
			"OidcController.GetDeviceTypeCookie",
			"Couldn't get device type cookie",
			err.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}

	var deviceIdValue *device_id.DeviceId
	deviceId, err := r.Cookie(constants.RoomPlayDeviceIdCookieName)
	if err != nil {
		deviceIdValue = nil
	} else {
		deviceIdValue = device_id.ParseDeviceId(deviceId.Value)
	}

	apiTokenResponse, err := handler.oidcAuthenticationService.AuthenticateWithGoogle(r.Context(), code, deviceIdValue, device_type.ParseDeviceType(&deviceType.Value))
	if err != nil {
		response.WriteJsonApplicationFailure(w, err, r.URL.RequestURI())
		return
	}
	base64RefreshToken := base64.RawURLEncoding.EncodeToString([]byte(apiTokenResponse.RefreshToken))
	cookie_helpers.SetAccessTokenCookie(w, apiTokenResponse.AccessToken)
	cookie_helpers.SetRefreshTokenCookie(w, base64RefreshToken)
	cookie_helpers.SetDeviceIdCookie(w, *apiTokenResponse.DeviceId.String())
	cookie_helpers.ClearStateCookie(w)

	http.Redirect(w, r, handler.configuration.Authentication().ClientOrigin+"/signin-oidc", http.StatusSeeOther)
}
