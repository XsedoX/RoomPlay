package oidc_authentication_service

import (
	"context"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_google_oidc_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/dtos/oidc_authenticate_user_service_dto"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
)

type OidcAuthenticationService struct {
	googleOidcService          i_google_oidc_service.IGoogleOidcService
	userRepository             i_user_repository.IUserRepository
	registerUserCommandHandler i_command_handler.ICommandHandlerWithResponse[*register_user_command.RegisterUserCommand, *register_user_command_response.RegisterUserCommandResponse]
	loginUserCommandHandler    i_command_handler.ICommandHandlerWithResponse[*login_user_command.LoginUserCommand, *login_user_command_response.LoginUserCommandResponse]
	unitOfWork                 i_unit_of_work.IUnitOfWork
}

func NewOidcAuthenticationService(googleOidcService i_google_oidc_service.IGoogleOidcService,
	userRepository i_user_repository.IUserRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	registerUserHandler i_command_handler.ICommandHandlerWithResponse[*register_user_command.RegisterUserCommand, *register_user_command_response.RegisterUserCommandResponse],
	loginUserHandler i_command_handler.ICommandHandlerWithResponse[*login_user_command.LoginUserCommand, *login_user_command_response.LoginUserCommandResponse],
) *OidcAuthenticationService {
	return &OidcAuthenticationService{
		googleOidcService:          googleOidcService,
		userRepository:             userRepository,
		unitOfWork:                 unitOfWork,
		registerUserCommandHandler: registerUserHandler,
		loginUserCommandHandler:    loginUserHandler,
	}
}

func (oidcAuthentication *OidcAuthenticationService) AuthenticateWithGoogle(ctx context.Context, code string, deviceId *device_id.DeviceId, deviceType *device_type.DeviceType) (*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto,
	error,
) {
	tokenResp, err := oidcAuthentication.googleOidcService.GetAccessToken(ctx, code)
	if err != nil {
		return nil, custom_error.NewCustomError("OidcAuthenticationService.GetAccessToken",
			"Couldn't get access token",
			err,
			custom_error_type.Unexpected)
	}

	claims, err := oidcAuthentication.googleOidcService.ParseIdToken(tokenResp.IdToken)
	if err != nil {
		return nil, custom_error.NewCustomError("OidcAuthenticationService.ParseIdToken",
			"Couldn't parse id token",
			err,
			custom_error_type.Unexpected)
	}
	var deviceTypeToPass device_type.DeviceType
	if deviceType == nil {
		deviceTypeToPass = device_type.Desktop
	} else {
		deviceTypeToPass = *deviceType
	}
	var apiTokenResponse oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer()) {
		// NOTE: User Login
		loginUserCommand := login_user_command.LoginUserCommand{
			Name: claims.GivenName,
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   deviceId,
				DeviceType: deviceTypeToPass,
			},
			Surname: claims.FamilyName,
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               claims.Subject,
				AccessToken:              tokenResp.AccessToken,
				RefreshToken:             tokenResp.RefreshToken,
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
			},
		}
		loginResponse, err := oidcAuthentication.loginUserCommandHandler.Handle(ctx, &loginUserCommand)
		if err != nil {
			return nil, err
		}
		apiTokenResponse.AccessToken = loginResponse.AccessToken
		apiTokenResponse.RefreshToken = loginResponse.RefreshToken
		apiTokenResponse.DeviceId = loginResponse.DeviceId
	} else {
		// NOTE: User Registration
		registerUserCommand := register_user_command.RegisterUserCommand{
			Name:       claims.GivenName,
			DeviceType: deviceTypeToPass,
			Surname:    claims.FamilyName,
			CredentialsDto: register_user_command.CredentialsDto{
				AccessToken:              tokenResp.AccessToken,
				ExternalId:               claims.Subject,
				RefreshToken:             tokenResp.RefreshToken,
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
			},
		}
		registerResponse, err := oidcAuthentication.registerUserCommandHandler.Handle(ctx, &registerUserCommand)
		if err != nil {
			return nil, err
		}
		apiTokenResponse.AccessToken = registerResponse.AccessToken
		apiTokenResponse.RefreshToken = registerResponse.RefreshToken
		apiTokenResponse.DeviceId = registerResponse.DeviceId
	}
	return &apiTokenResponse, err
}
