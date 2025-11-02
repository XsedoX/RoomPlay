package login

import (
	"context"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type UserCommandHandler struct {
	unitOfWork                    contracts.IUnitOfWork
	userRepository                contracts.IUserRepository
	externalCredentialsRepository contracts.IExternalCredentialsRepository
	encrypter                     contracts.IEncrypter
	jwtProvider                   contracts.IJwtProvider
	refreshTokenRepository        contracts.IRefreshTokenRepository
}

func NewLoginUserCommandHandler(unitOfWork contracts.IUnitOfWork,
	userRepository contracts.IUserRepository,
	encrypter contracts.IEncrypter,
	jwtProvider contracts.IJwtProvider,
	refreshTokenRepository contracts.IRefreshTokenRepository,
	externalCredentialsRepository contracts.IExternalCredentialsRepository) *UserCommandHandler {
	return &UserCommandHandler{
		unitOfWork:                    unitOfWork,
		userRepository:                userRepository,
		encrypter:                     encrypter,
		jwtProvider:                   jwtProvider,
		refreshTokenRepository:        refreshTokenRepository,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (handler *UserCommandHandler) Handle(ctx context.Context, command *UserCommand) (*UserCommandResponse, error) {
	var response UserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		userFromDb, err := handler.userRepository.GetUserByExternalId(ctx, command.ExternalId, handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}

		userFromDb.ChangeFullName(*user.NewFullName(command.Name, command.Surname))

		var deviceId shared.DeviceId
		if command.DeviceDto.DeviceId == nil {
			deviceId = userFromDb.LoginWithNewDevice(command.DeviceDto.DeviceType)
		} else {
			userFromDb.ReloginWithKnownDevice(*command.DeviceDto.DeviceId)
			deviceId = *command.DeviceDto.DeviceId
		}
		response.DeviceId = deviceId
		err = handler.userRepository.Update(ctx, userFromDb, handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userFromDb.Id(), deviceId, string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return newTokenErr
		}
		response.RefreshToken = string(refreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return tokenErr
		}

		creds := credentials.NewExternalCredentials(userFromDb.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.Scopes,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		return handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
