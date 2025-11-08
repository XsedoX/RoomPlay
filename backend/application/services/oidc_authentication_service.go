package services

import (
	"context"
	"time"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/dtos"
	"xsedox.com/main/application/user/login_command"
	"xsedox.com/main/application/user/register_command"
	"xsedox.com/main/domain/user"
)

type OidcAuthenticationService struct {
	googleOidcService          contracts.IGoogleOidcService
	userRepository             contracts.IUserRepository
	registerUserCommandHandler contracts.ICommandHandlerWithResponse[*register.UserCommand, *register.UserCommandResponse]
	loginUserCommandHandler    contracts.ICommandHandlerWithResponse[*login_command.UserCommand, *login_command.UserCommandResponse]
	unitOfWork                 contracts.IUnitOfWork
}

func NewOidcAuthenticationService(googleOidcService contracts.IGoogleOidcService,
	userRepository contracts.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	registerUserHandler contracts.ICommandHandlerWithResponse[*register.UserCommand, *register.UserCommandResponse],
	loginUserHandler contracts.ICommandHandlerWithResponse[*login_command.UserCommand, *login_command.UserCommandResponse]) *OidcAuthenticationService {
	return &OidcAuthenticationService{
		googleOidcService:          googleOidcService,
		userRepository:             userRepository,
		unitOfWork:                 unitOfWork,
		registerUserCommandHandler: registerUserHandler,
		loginUserCommandHandler:    loginUserHandler,
	}
}

func (oidcAuthentication *OidcAuthenticationService) AuthenticateWithGoogle(ctx context.Context, code string, deviceId *user.DeviceId, deviceType *user.DeviceType) (*dtos.OidcAuthenticateUserServiceDto,
	error) {
	tokenResp, err := oidcAuthentication.googleOidcService.GetAccessToken(ctx, code)
	if err != nil {
		return nil, applicationErrors.NewApplicationError("OidcAuthenticationService.GetAccessToken",
			"Couldn't get access token",
			err,
			applicationErrors.Unexpected)
	}

	claims, err := oidcAuthentication.googleOidcService.ParseIdToken(tokenResp.IdToken)
	if err != nil {
		return nil, applicationErrors.NewApplicationError("OidcAuthenticationService.ParseIdToken",
			"Couldn't parse id token",
			err,
			applicationErrors.Unexpected)
	}
	var deviceTypeToPass user.DeviceType
	if deviceType == nil {
		deviceTypeToPass = user.Desktop
	} else {
		deviceTypeToPass = *deviceType
	}
	var apiTokenResponse dtos.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer()) {
		loginUserCommand := login_command.UserCommand{
			Name: claims.GivenName,
			DeviceDto: login_command.DeviceDto{
				DeviceId:   deviceId,
				DeviceType: deviceTypeToPass,
			},
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: login_command.CredentialsDto{
				AccessToken:              tokenResp.AccessToken,
				RefreshToken:             tokenResp.RefreshToken,
				Scopes:                   tokenResp.Scope,
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
		registerUserCommand := register.UserCommand{
			Name:       claims.GivenName,
			DeviceType: deviceTypeToPass,
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: register.CredentialsDto{
				AccessToken:              tokenResp.AccessToken,
				RefreshToken:             tokenResp.RefreshToken,
				Scopes:                   tokenResp.Scope,
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
