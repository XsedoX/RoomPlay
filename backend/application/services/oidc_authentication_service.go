package services

import (
	"context"
	"time"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/dtos"
	"xsedox.com/main/application/user/login"
	"xsedox.com/main/application/user/register"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/shared"
)

type OidcAuthenticationService struct {
	googleOidcService   contracts.IGoogleOidcService
	userRepository      contracts.IUserRepository
	registerUserHandler contracts.ICommandHandlerWithResponse[*register.UserCommand, *register.UserCommandResponse]
	loginUserHandler    contracts.ICommandHandlerWithResponse[*login.UserCommand, *login.UserCommandResponse]
	unitOfWork          contracts.IUnitOfWork
}

func NewOidcAuthenticationService(googleOidcService contracts.IGoogleOidcService,
	userRepository contracts.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	registerUserHandler contracts.ICommandHandlerWithResponse[*register.UserCommand, *register.UserCommandResponse],
	loginUserHandler contracts.ICommandHandlerWithResponse[*login.UserCommand, *login.UserCommandResponse]) *OidcAuthenticationService {
	return &OidcAuthenticationService{
		googleOidcService:   googleOidcService,
		userRepository:      userRepository,
		unitOfWork:          unitOfWork,
		registerUserHandler: registerUserHandler,
		loginUserHandler:    loginUserHandler,
	}
}

func (oidcAuthentication *OidcAuthenticationService) AuthenticateWithGoogle(ctx context.Context, code string, deviceId *shared.DeviceId) (*dtos.OidcAuthenticateUserServiceDto,
	*applicationErrors.ApplicationError) {
	tokenResp, err := oidcAuthentication.googleOidcService.GetAccessToken(ctx, code)
	if err != nil {
		return nil, applicationErrors.NewApplicationError("couldn't get access token", err,
			applicationErrors.ErrInfrastructure)
	}

	claims, err := oidcAuthentication.googleOidcService.ParseIdToken(tokenResp.AccessToken)
	if err != nil {
		return nil, applicationErrors.NewApplicationError("Couldn't parse id token", err,
			applicationErrors.ErrInfrastructure)

	}

	var apiTokenResponse dtos.OidcAuthenticateUserServiceDto
	if oidcAuthentication.userRepository.CheckIfUserExistByExternalId(ctx, claims.Subject, oidcAuthentication.unitOfWork.GetQueryer()) {
		loginUserCommand := login.UserCommand{
			Name: claims.GivenName,
			DeviceDto: login.DeviceDto{
				DeviceId:   deviceId,
				DeviceType: device.Computer,
			},
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: login.CredentialsDto{
				AccessToken:           tokenResp.AccessToken,
				RefreshToken:          tokenResp.RefreshToken,
				Scopes:                tokenResp.Scope,
				AccessTokenExpiresAt:  time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
				RefreshTokenExpiresAt: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
			},
		}
		loginResponse, err := oidcAuthentication.loginUserHandler.Handle(ctx, &loginUserCommand)
		if err != nil {
			return nil, applicationErrors.NewApplicationError("couldn't create login user command", err, applicationErrors.ErrInfrastructure)
		}
		apiTokenResponse.AccessToken = loginResponse.AccessToken
		apiTokenResponse.RefreshToken = loginResponse.RefreshToken
		apiTokenResponse.DeviceId = loginResponse.DeviceId
	} else {
		registerUserCommand := register.UserCommand{
			Name:       claims.GivenName,
			DeviceType: device.Computer,
			ExternalId: claims.Subject,
			Surname:    claims.FamilyName,
			CredentialsDto: register.CredentialsDto{
				AccessToken:           tokenResp.AccessToken,
				RefreshToken:          tokenResp.RefreshToken,
				Scopes:                tokenResp.Scope,
				AccessTokenExpiresAt:  time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
				RefreshTokenExpiresAt: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UTC(),
			},
		}
		registerResponse, err := oidcAuthentication.registerUserHandler.Handle(ctx, &registerUserCommand)
		if err != nil {
			return nil, applicationErrors.NewApplicationError("couldn't create register user command", err, applicationErrors.ErrInfrastructure)
		}
		apiTokenResponse.AccessToken = registerResponse.AccessToken
		apiTokenResponse.RefreshToken = registerResponse.RefreshToken
		apiTokenResponse.DeviceId = registerResponse.DeviceId
	}

	return &apiTokenResponse, nil
}
