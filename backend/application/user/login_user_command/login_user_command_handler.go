package login_user_command

import (
	"context"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts2 "xsedox.com/main/application/user/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
)

type LoginUserCommandHandler struct {
	unitOfWork                    contracts.IUnitOfWork
	userRepository                contracts2.IUserRepository
	externalCredentialsRepository contracts.IExternalCredentialsRepository
	encrypter                     contracts.IEncrypter
	jwtProvider                   contracts.IJwtProvider
	refreshTokenRepository        contracts.IRefreshTokenRepository
}

func NewLoginUserCommandHandler(unitOfWork contracts.IUnitOfWork,
	userRepository contracts2.IUserRepository,
	encrypter contracts.IEncrypter,
	jwtProvider contracts.IJwtProvider,
	refreshTokenRepository contracts.IRefreshTokenRepository,
	externalCredentialsRepository contracts.IExternalCredentialsRepository) *LoginUserCommandHandler {
	return &LoginUserCommandHandler{
		unitOfWork:                    unitOfWork,
		userRepository:                userRepository,
		encrypter:                     encrypter,
		jwtProvider:                   jwtProvider,
		refreshTokenRepository:        refreshTokenRepository,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (handler *LoginUserCommandHandler) Handle(ctx context.Context, command *LoginUserCommand) (*LoginUserCommandResponse, error) {
	var response LoginUserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		userFromDb, err := handler.userRepository.GetUserByExternalId(ctx, command.ExternalId, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("LoginUserCommandHandler.GetUserByExternalId",
				"Problem with getting user with external id",
				err,
				custom_errors.Unexpected)
		}

		userFromDb.ChangeFullName(user.NewFullName(command.Name, command.Surname))

		var deviceId user.DeviceId
		if command.DeviceDto.DeviceId == nil {
			deviceId = userFromDb.LoginWithNewDevice(command.DeviceDto.DeviceType)
		} else {
			userFromDb.ReloginWithKnownDevice(*command.DeviceDto.DeviceId)
			deviceId = *command.DeviceDto.DeviceId
		}
		response.DeviceId = deviceId
		err = handler.userRepository.Update(ctx, userFromDb, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("LoginUserCommandHandler.Update",
				"Problem with updating user in the database",
				err,
				custom_errors.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userFromDb.Id(), deviceId, string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return custom_errors.NewCustomError("LoginUserCommandHandler.AssignNewToken",
				"Problem with assigning new token to user",
				newTokenErr,
				custom_errors.Unexpected)
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
		grantErr := handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
		if grantErr != nil {
			return custom_errors.NewCustomError("LoginUserCommandHandler.Grant",
				"Problem with assigning external credentials to user",
				grantErr,
				custom_errors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
