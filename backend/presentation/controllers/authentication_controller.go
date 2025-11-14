package controllers

import (
	"encoding/base64"
	"net/http"

	"xsedox.com/main/application"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/user/login_user_refresh_token_command"
	"xsedox.com/main/application/user/logout_user_command"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/presentation/helpers"
	"xsedox.com/main/presentation/response"
)

type AuthenticationController struct {
	loginRefreshTokenCommandHandler  contracts.ICommandHandlerWithResponse[*string, *login_user_refresh_token_command.LoginUserRefreshTokenCommandResponse]
	configuration                    config.IConfiguration
	logoutRefreshTokenCommandHandler contracts.ICommandHandler[*logout_user_command.LogoutUserCommand]
}

func NewAuthenticationController(refreshTokenCommandHandler contracts.ICommandHandlerWithResponse[*string, *login_user_refresh_token_command.LoginUserRefreshTokenCommandResponse],
	configuration config.IConfiguration,
	logoutRefreshTokenCommandHandler contracts.ICommandHandler[*logout_user_command.LogoutUserCommand]) *AuthenticationController {
	return &AuthenticationController{
		loginRefreshTokenCommandHandler:  refreshTokenCommandHandler,
		configuration:                    configuration,
		logoutRefreshTokenCommandHandler: logoutRefreshTokenCommandHandler,
	}
}

func (handler *AuthenticationController) RefreshToken(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := req.Cookie(helpers.RoomplayRefreshTokenCookieName)
	if err != nil || refreshToken == nil || refreshToken.Value == "" {
		helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
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
		helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
		return
	}
	encodedRefreshToken := base64.RawURLEncoding.EncodeToString([]byte(result.RefreshToken))
	helpers.SetRefreshTokenCookie(w, encodedRefreshToken, handler.configuration.Server().BasePath)
	helpers.SetAccessTokenCookie(w, result.AccessToken, handler.configuration.Server().BasePath)
	response.WriteJsonNoContent(w)
}

func (handler *AuthenticationController) Logout(w http.ResponseWriter, req *http.Request) {
	var command logout_user_command.LogoutUserCommand
	userId, ok := application.GetUserIdFromContext(req.Context())
	if !ok {
		response.WriteJsonApplicationFailure(w,
			application.NewMissingUserIdInContextError,
			req.URL.RequestURI(),
		)
		helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		return
	}
	command.UserId = *userId

	deviceId, err := req.Cookie(helpers.RoomPlayDeviceIdCookieName)
	if deviceId == nil || err != nil {
		helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonFailure(w,
			"AuthenticationController.RoomPlayDeviceIdCookie",
			"Missing device id cookie",
			err.Error(),
			req.URL.RequestURI(),
			http.StatusPartialContent)
		command.DeviceId = nil
	} else {
		deviceIdValue := deviceId.Value
		command.DeviceId = user.ParseDeviceId(deviceIdValue)
		response.WriteJsonNoContent(w)
	}
	err = handler.logoutRefreshTokenCommandHandler.Handle(req.Context(), &command)
	if err != nil {
		helpers.ClearRefreshTokenCookie(w, handler.configuration.Server().BasePath)
		helpers.ClearAccessTokenCookie(w, handler.configuration.Server().BasePath)
		response.WriteJsonApplicationFailure(w,
			err,
			req.URL.RequestURI(),
		)
	}
}
