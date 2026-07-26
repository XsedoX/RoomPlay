package external_authentication_service

import (
	"context"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_authentication_service_provider"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/dtos/oidc_authenticate_user_service_dto"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type ExternalAuthenticationService struct {
	authServices                  map[music_provider.MusicProvider]i_external_authentication_service_provider.IExternalAuthenticationServiceProvider
	userRepository                i_user_repository.IUserRepository
	registerUserCommandHandler    i_command_handler.ICommandHandlerWithResponse[*register_user_command.RegisterUserCommand, *register_user_command_response.RegisterUserCommandResponse]
	loginUserCommandHandler       i_command_handler.ICommandHandlerWithResponse[*login_user_command.LoginUserCommand, *login_user_command_response.LoginUserCommandResponse]
	unitOfWork                    i_unit_of_work.IUnitOfWork
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
}

func NewExternalAuthenticationService(
	authServices map[music_provider.MusicProvider]i_external_authentication_service_provider.IExternalAuthenticationServiceProvider,
	userRepository i_user_repository.IUserRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	registerUserHandler i_command_handler.ICommandHandlerWithResponse[*register_user_command.RegisterUserCommand, *register_user_command_response.RegisterUserCommandResponse],
	loginUserHandler i_command_handler.ICommandHandlerWithResponse[*login_user_command.LoginUserCommand, *login_user_command_response.LoginUserCommandResponse],
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository,
) *ExternalAuthenticationService {
	return &ExternalAuthenticationService{
		authServices:                  authServices,
		userRepository:                userRepository,
		unitOfWork:                    unitOfWork,
		registerUserCommandHandler:    registerUserHandler,
		loginUserCommandHandler:       loginUserHandler,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (oidcAuthentication *ExternalAuthenticationService) AuthenticateWithExternalProvider(
	ctx context.Context,
	code string,
	deviceId *device_id.DeviceId,
	deviceType *device_type.DeviceType,
	provider music_provider.MusicProvider,
) (*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto,
	error,
) {
	oidcService, ok := oidcAuthentication.authServices[provider]
	if !ok {
		return nil, application_error.NewApplicationError("OidcAuthenticationService.AuthenticateWithExternalProvider.MissingProviderInDI",
			"Missing provider in DI",
			nil,
			application_error_type.Unexpected,
		)
	}
	tokenResp, err := oidcService.GetAccessToken(ctx, code)
	if err != nil {
		return nil, application_error.NewApplicationError("OidcAuthenticationService.GetAccessToken",
			"Couldn't get access token",
			err,
			application_error_type.Unexpected)
	}

	claims, err := oidcService.ParseIdToken(tokenResp.IdToken)
	if err != nil {
		return nil, application_error.NewApplicationError("OidcAuthenticationService.ParseIdToken",
			"Couldn't parse id token",
			err,
			application_error_type.Unexpected)
	}
	var deviceTypeToPass device_type.DeviceType
	if deviceType == nil {
		deviceTypeToPass = device_type.Desktop
	} else {
		deviceTypeToPass = *deviceType
	}
	var apiTokenResponse oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer(ctx)) {
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

func (oidcAuthentication *ExternalAuthenticationService) RefreshAccessTokenWithExternalProvider(ctx context.Context, userId user_id.UserId) (*token.Token, error) {
	externalCreds, err := oidcAuthentication.externalCredentialsRepository.GetExternalCredentialsByUserId(ctx, userId, oidcAuthentication.unitOfWork.GetQueryer(ctx))
	if err != nil {
		return nil, application_error.NewApplicationError("ExternalAuthenticationService.RefreshTokenWithExternalProvider.GetRefreshTokenByUserId",
			"Problem with getting refresh token for music service.",
			err,
			application_error_type.Unexpected,
		)
	}
	oidcService, ok := oidcAuthentication.authServices[externalCreds.MusicProvider()]
	if !ok {
		return nil, application_error.NewApplicationError("ExternalAuthenticationService.RefreshTokenWithExternalProvider.MissingProviderInDI",
			"Missing provider in DI",
			nil,
			application_error_type.Unexpected,
		)
	}
	if externalCreds.GetRefreshToken().IsExpired() {
		return nil,
			application_error.NewApplicationError("ExternalAuthenticationService.RefreshAccessTokenWithExternalProvider.RefreshTokenExpired",
				"Refresh token is expired, user needs to re-authenticate.",
				nil,
				application_error_type.Unauthorized,
			)
	}
	newAccessTokenDto, err := oidcService.RefreshAccessToken(
		ctx,
		externalCreds.GetRefreshToken().Value(),
	)
	if err != nil {
		return nil, application_error.NewApplicationError("ExternalAuthenticationService.RefreshAccessTokenWithExternalProvider.RefreshAccessToken",
			"Problem with refreshing access token with external provider.",
			err,
			application_error_type.Unexpected,
		)
	}
	newAccessToken, err := token.NewToken(newAccessTokenDto.AccessToken, time.Now().Add(time.Duration(newAccessTokenDto.ExpiresIn)*time.Second).UTC())
	if err != nil {
		return nil, err
	}
	externalCreds.UpdateAccessToken(*newAccessToken)
	updateErr := oidcAuthentication.externalCredentialsRepository.UpdateExternalCredentials(ctx, externalCreds, oidcAuthentication.unitOfWork.GetQueryer(ctx))
	if updateErr != nil {
		return nil, application_error.NewApplicationError("ExternalAuthenticationService.RefreshAccessTokenWithExternalProvider.UpdateExternalCredentials",
			"Problem with updating external credentials in the database.",
			updateErr,
			application_error_type.Unexpected,
		)
	}
	return newAccessToken, nil
}
