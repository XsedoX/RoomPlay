package register_user

import (
	"context"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts2 "xsedox.com/main/application/user/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
)

type RegisterUserCommandHandler struct {
	userRepository                contracts2.IUserRepository
	externalCredentialsRepository contracts.IExternalCredentialsRepository
	refreshTokenRepository        contracts.IRefreshTokenRepository
	unitOfWork                    contracts.IUnitOfWork
	jwtProvider                   contracts.IJwtProvider
	encrypter                     contracts.IEncrypter
}

func NewRegisterUserCommandHandler(userRepository contracts2.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	credsRepository contracts.IExternalCredentialsRepository,
	jwtProvider contracts.IJwtProvider,
	refreshTokenRepository contracts.IRefreshTokenRepository,
	encrypter contracts.IEncrypter) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{
		userRepository:                userRepository,
		externalCredentialsRepository: credsRepository,
		unitOfWork:                    unitOfWork,
		jwtProvider:                   jwtProvider,
		refreshTokenRepository:        refreshTokenRepository,
		encrypter:                     encrypter,
	}
}
func (handler *RegisterUserCommandHandler) Handle(ctx context.Context, command *RegisterUserCommand) (*RegisterUserCommandResponse, error) {
	var response RegisterUserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		userAgg := user.NewUser(command.ExternalId, command.Name, command.Surname, command.DeviceType)
		deviceEnt := userAgg.GetMostRecentDevice()
		err := handler.userRepository.Add(ctx, userAgg, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("RegisterUserCommandHandler.UserRepository.Add",
				"Adding user problem",
				err,
				custom_errors.Unexpected)
		}
		response.DeviceId = deviceEnt.Id()

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userAgg.Id())
		if tokenErr != nil {
			return custom_errors.NewCustomError("RegisterUserCommandHandler.GenerateToken",
				"Access token generation problem",
				tokenErr,
				custom_errors.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userAgg.Id(), deviceEnt.Id(), string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return custom_errors.NewCustomError("RegisterUserCommandHandler.AssignNewToken",
				"Access token generation problem",
				newTokenErr,
				custom_errors.Unexpected)
		}
		response.RefreshToken = string(refreshTokenValue)

		creds := credentials.NewExternalCredentials(userAgg.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.Scopes,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		grantErr := handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
		if grantErr != nil {
			return custom_errors.NewCustomError("RegisterUserCommandHandler.Grant",
				"Problem with granting external credentials.",
				grantErr,
				custom_errors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, err
}
