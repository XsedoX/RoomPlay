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

func (handler *AuthenticationController) RefreshToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := req.Cookie(constants.RoomPlayRefreshTokenCookieName)
	if err != nil || refreshToken == nil || refreshToken.Value == "" {
		cookie_helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		cookie_helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonFailure(w,
			"AuthenticationController.MissingRefreshTokenCookie",
			"Missing refresh token cookie",
			"Cookie with refresh token had issues.",
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
		cookie_helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		cookie_helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
		return
	}
	encodedRefreshToken := base64.RawURLEncoding.EncodeToString([]byte(result.RefreshToken))
	cookie_helpers.SetRefreshTokenCookie(w, encodedRefreshToken, handler.configuration.Server().BasePath)
	cookie_helpers.SetAccessTokenCookie(w, result.AccessToken, handler.configuration.Server().BasePath)
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
		cookie_helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		cookie_helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		return
	}
	command.UserId = *userId

	deviceId, err := req.Cookie(constants.RoomPlayDeviceIdCookieName)
	if deviceId == nil || err != nil {
		cookie_helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		cookie_helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
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
		cookie_helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		cookie_helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
	}
}
