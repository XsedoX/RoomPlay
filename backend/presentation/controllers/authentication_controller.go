package controllers

import (
	"encoding/base64"
	"net/http"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token"
	"github.com/XsedoX/RoomPlay/application/user/logout_user"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/presentation/helpers"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

const (
	AuthBasePath   = "/auth"
	LogoutBasePath = "/logout"
)

type AuthenticationController struct {
	loginRefreshTokenCommandHandler  contracts.ICommandHandlerWithResponse[*string, *login_user_refresh_token.LoginUserRefreshTokenCommandResponse]
	configuration                    config.IConfiguration
	logoutRefreshTokenCommandHandler contracts.ICommandHandler[*logout_user.LogoutUserCommand]
}

func NewAuthenticationController(refreshTokenCommandHandler contracts.ICommandHandlerWithResponse[*string, *login_user_refresh_token.LoginUserRefreshTokenCommandResponse],
	configuration config.IConfiguration,
	logoutRefreshTokenCommandHandler contracts.ICommandHandler[*logout_user.LogoutUserCommand],
) *AuthenticationController {
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
	var command logout_user.LogoutUserCommand
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
