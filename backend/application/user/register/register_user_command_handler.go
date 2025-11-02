package register

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/user"
)

type UserCommandHandler struct {
	userRepository                contracts.IUserRepository
	externalCredentialsRepository contracts.IExternalCredentialsRepository
	refreshTokenRepository        contracts.IRefreshTokenRepository
	unitOfWork                    contracts.IUnitOfWork
	jwtProvider                   contracts.IJwtProvider
	encrypter                     contracts.IEncrypter
}

func NewRegisterUserCommandHandler(userRepository contracts.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	credsRepository contracts.IExternalCredentialsRepository,
	jwtProvider contracts.IJwtProvider,
	refreshTokenRepository contracts.IRefreshTokenRepository,
	encrypter contracts.IEncrypter) *UserCommandHandler {
	return &UserCommandHandler{
		userRepository:                userRepository,
		externalCredentialsRepository: credsRepository,
		unitOfWork:                    unitOfWork,
		jwtProvider:                   jwtProvider,
		refreshTokenRepository:        refreshTokenRepository,
		encrypter:                     encrypter,
	}
}
func (handler *UserCommandHandler) Handle(ctx context.Context, command *UserCommand) (*UserCommandResponse, error) {
	var response UserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		deviceEnt := device.NewDevice(command.DeviceType)
		userAgg := user.NewUser(command.ExternalId, command.Name, command.Surname, *deviceEnt)
		err := handler.userRepository.Add(ctx, userAgg, handler.unitOfWork.GetQueryer())
		response.DeviceId = deviceEnt.Id()
		if err != nil {
			return err
		}

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userAgg.Id())
		if tokenErr != nil {
			return tokenErr
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userAgg.Id(), deviceEnt.Id(), string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return newTokenErr
		}
		response.RefreshToken = string(refreshTokenValue)

		creds := credentials.NewExternalCredentials(userAgg.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.Scopes,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		return handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
	})
	if err != nil {
		respErr := applicationErrors.NewApplicationError("problem with executing transaction", err, applicationErrors.Unexpected)
		return nil, respErr
	}
	return &response, nil
}
