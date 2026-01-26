package services

import (
	"context"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/dtos"
	"github.com/XsedoX/RoomPlay/application/user/login_user"
	"github.com/XsedoX/RoomPlay/application/user/register_user"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type OidcAuthenticationService struct {
	googleOidcService          application_contracts.IGoogleOidcService
	userRepository             user_contracts.IUserRepository
	registerUserCommandHandler application_contracts.ICommandHandlerWithResponse[*register_user.RegisterUserCommand, *register_user.RegisterUserCommandResponse]
	loginUserCommandHandler    application_contracts.ICommandHandlerWithResponse[*login_user.LoginUserCommand, *login_user.LoginUserCommandResponse]
	unitOfWork                 application_contracts.IUnitOfWork
}

func NewOidcAuthenticationService(googleOidcService application_contracts.IGoogleOidcService,
	userRepository user_contracts.IUserRepository,
	unitOfWork application_contracts.IUnitOfWork,
	registerUserHandler application_contracts.ICommandHandlerWithResponse[*register_user.RegisterUserCommand, *register_user.RegisterUserCommandResponse],
	loginUserHandler application_contracts.ICommandHandlerWithResponse[*login_user.LoginUserCommand, *login_user.LoginUserCommandResponse],
) *OidcAuthenticationService {
	return &OidcAuthenticationService{
		googleOidcService:          googleOidcService,
		userRepository:             userRepository,
		unitOfWork:                 unitOfWork,
		registerUserCommandHandler: registerUserHandler,
		loginUserCommandHandler:    loginUserHandler,
	}
}

func (oidcAuthentication *OidcAuthenticationService) AuthenticateWithGoogle(ctx context.Context, code string, deviceId *user.DeviceId, deviceType *user.DeviceType) (*dtos.OidcAuthenticateUserServiceDto,
	error,
) {
	tokenResp, err := oidcAuthentication.googleOidcService.GetAccessToken(ctx, code)
	if err != nil {
		return nil, customerrors.NewCustomError("OidcAuthenticationService.GetAccessToken",
			"Couldn't get access token",
			err,
			customerrors.Unexpected)
	}

	claims, err := oidcAuthentication.googleOidcService.ParseIdToken(tokenResp.IdToken)
	if err != nil {
		return nil, customerrors.NewCustomError("OidcAuthenticationService.ParseIdToken",
			"Couldn't parse id token",
			err,
			customerrors.Unexpected)
	}
	var deviceTypeToPass user.DeviceType
	if deviceType == nil {
		deviceTypeToPass = user.Desktop
	} else {
		deviceTypeToPass = *deviceType
	}
	var apiTokenResponse dtos.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer()) {
		// NOTE: User Login
		loginUserCommand := login_user.LoginUserCommand{
			Name: claims.GivenName,
			DeviceDto: login_user.DeviceDto{
				DeviceId:   deviceId,
				DeviceType: deviceTypeToPass,
			},
			Surname: claims.FamilyName,
			CredentialsDto: login_user.CredentialsDto{
				ExternalId:               claims.Subject,
				AccessToken:              tokenResp.AccessToken,
				RefreshToken:             tokenResp.RefreshToken,
				MusicProvider:            credentials.YouTube,
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
		registerUserCommand := register_user.RegisterUserCommand{
			Name:       claims.GivenName,
			DeviceType: deviceTypeToPass,
			Surname:    claims.FamilyName,
			CredentialsDto: register_user.CredentialsDto{
				AccessToken:              tokenResp.AccessToken,
				ExternalId:               claims.Subject,
				RefreshToken:             tokenResp.RefreshToken,
				MusicProvider:            credentials.YouTube,
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
