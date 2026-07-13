package authentication_controller

import (
	"encoding/base64"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token/login_user_refresh_token_command_response"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/cookie_helpers"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

type AuthenticationController struct {
	loginRefreshTokenCommandHandler  i_command_handler.ICommandHandlerWithResponse[*string, *login_user_refresh_token_command_response.LoginUserRefreshTokenCommandResponse]
	configuration                    config.IConfiguration
	logoutRefreshTokenCommandHandler i_command_handler.ICommandHandler[*logout_user_command.LogoutUserCommand]
}

func NewAuthenticationController(refreshTokenCommandHandler i_command_handler.ICommandHandlerWithResponse[*string, *login_user_refresh_token_command_response.LoginUserRefreshTokenCommandResponse],
	configuration config.IConfiguration,
	logoutRefreshTokenCommandHandler i_command_handler.ICommandHandler[*logout_user_command.LogoutUserCommand],
) *AuthenticationController {
	return &AuthenticationController{
		loginRefreshTokenCommandHandler:  refreshTokenCommandHandler,
		configuration:                    configuration,
		logoutRefreshTokenCommandHandler: logoutRefreshTokenCommandHandler,
	}
}

// RefreshToken handles refreshing the user's access token using a refresh token stored in a cookie. It validates the refresh token, generates new tokens, and updates the cookies accordingly.
// @Summary Refresh user access token
// @Description Refreshes the user's access token using a refresh token stored in a cookie. Validates the refresh token, generates new tokens, and updates the cookies accordingly.
// @Tags auth
// @Accept json
// @Produce json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router /auth/refresh-token [post]
// @Security BearerAuth
func (handler *AuthenticationController) RefreshToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := req.Cookie(constants.RoomPlayRefreshTokenCookieName)
	if err != nil || refreshToken == nil || refreshToken.Value == "" {
		cookie_helpers.ClearRefreshTokenCookie(w)
		cookie_helpers.ClearAccessTokenCookie(w)
		response.WriteJsonFailure(w,
			"AuthenticationController.MissingRefreshTokenCookie",
			"Missing refresh token cookie",
			"Something went wrong while trying to reauthenticate you. Please log in again.",
			req.URL.RequestURI(),
			http.StatusUnauthorized)
		return
	}

	decodedToken, err := base64.RawURLEncoding.DecodeString(refreshToken.Value)
	if err != nil {
		// Handle malformed cookie value.
		response.WriteJsonFailure(w,
			"AuthenticationController.DecodeString",
			"Invalid refresh token",
			err.Error(),
			req.URL.RequestURI(),
			http.StatusUnauthorized)
		return
	}
	decodedTokenString := string(decodedToken)

	result, err := handler.loginRefreshTokenCommandHandler.Handle(req.Context(), &decodedTokenString)
	if err != nil {
		cookie_helpers.ClearRefreshTokenCookie(w)
		cookie_helpers.ClearAccessTokenCookie(w)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
		return
	}
	encodedRefreshToken := base64.RawURLEncoding.EncodeToString([]byte(result.RefreshToken))
	cookie_helpers.SetRefreshTokenCookie(w, encodedRefreshToken)
	cookie_helpers.SetAccessTokenCookie(w, result.AccessToken)
	response.WriteJsonNoContent(w)
}

func (handler *AuthenticationController) Logout(w http.ResponseWriter, req *http.Request) {
	var command logout_user_command.LogoutUserCommand
	userId, ok := application_helpers.GetUserIdFromContext(req.Context())
	if !ok {
		response.WriteJsonApplicationFailure(w,
			application_helpers.NewMissingUserIdInContextError,
			req.URL.RequestURI(),
		)
		cookie_helpers.ClearRefreshTokenCookie(w)
		cookie_helpers.ClearAccessTokenCookie(w)
		return
	}
	command.UserId = *userId

	deviceId, err := req.Cookie(constants.RoomPlayDeviceIdCookieName)
	if deviceId == nil || err != nil {
		cookie_helpers.ClearRefreshTokenCookie(w)
		cookie_helpers.ClearAccessTokenCookie(w)
		response.WriteJsonFailure(w,
			"AuthenticationController.RoomPlayDeviceIdCookie",
			"Missing device id cookie",
			err.Error(),
			req.URL.RequestURI(),
			http.StatusPartialContent)
		command.DeviceId = nil
	} else {
		deviceIdValue := deviceId.Value
		command.DeviceId = device_id.ParseDeviceId(deviceIdValue)
		response.WriteJsonNoContent(w)
	}
	err = handler.logoutRefreshTokenCommandHandler.Handle(req.Context(), &command)
	if err != nil {
		cookie_helpers.ClearRefreshTokenCookie(w)
		cookie_helpers.ClearAccessTokenCookie(w)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
	}
}
