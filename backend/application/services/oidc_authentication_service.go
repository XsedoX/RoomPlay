package services

import (
	"context"
	"time"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/application/dtos"
	contracts2 "xsedox.com/main/application/user/contracts"
	"xsedox.com/main/application/user/login_user"
	"xsedox.com/main/application/user/register_user"
	"xsedox.com/main/domain/user"
)

type OidcAuthenticationService struct {
	googleOidcService          contracts.IGoogleOidcService
	userRepository             contracts2.IUserRepository
	registerUserCommandHandler contracts.ICommandHandlerWithResponse[*register_user.RegisterUserCommand, *register_user.RegisterUserCommandResponse]
	loginUserCommandHandler    contracts.ICommandHandlerWithResponse[*login_user.LoginUserCommand, *login_user.LoginUserCommandResponse]
	unitOfWork                 contracts.IUnitOfWork
}

func NewOidcAuthenticationService(googleOidcService contracts.IGoogleOidcService,
	userRepository contracts2.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	registerUserHandler contracts.ICommandHandlerWithResponse[*register_user.RegisterUserCommand, *register_user.RegisterUserCommandResponse],
	loginUserHandler contracts.ICommandHandlerWithResponse[*login_user.LoginUserCommand, *login_user.LoginUserCommandResponse]) *OidcAuthenticationService {
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
		return nil, custom_errors.NewCustomError("OidcAuthenticationService.GetAccessToken",
			"Couldn't get access token",
			err,
			custom_errors.Unexpected)
	}

	claims, err := oidcAuthentication.googleOidcService.ParseIdToken(tokenResp.IdToken)
	if err != nil {
		return nil, custom_errors.NewCustomError("OidcAuthenticationService.ParseIdToken",
			"Couldn't parse id token",
			err,
			custom_errors.Unexpected)
	}
	var deviceTypeToPass user.DeviceType
	if deviceType == nil {
		deviceTypeToPass = user.Desktop
	} else {
		deviceTypeToPass = *deviceType
	}
	var apiTokenResponse dtos.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer()) {
		loginUserCommand := login_user.LoginUserCommand{
			Name: claims.GivenName,
			DeviceDto: login_user.DeviceDto{
				DeviceId:   deviceId,
				DeviceType: deviceTypeToPass,
			},
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: login_user.CredentialsDto{
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
		registerUserCommand := register_user.RegisterUserCommand{
			Name:       claims.GivenName,
			DeviceType: deviceTypeToPass,
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: register_user.CredentialsDto{
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
