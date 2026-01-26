package login_user

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type LoginUserCommandHandler struct {
	unitOfWork                    application_contracts.IUnitOfWork
	userRepository                user_contracts.IUserRepository
	externalCredentialsRepository application_contracts.IExternalCredentialsRepository
	encrypter                     application_contracts.IEncrypter
	jwtProvider                   application_contracts.IJwtProvider
	refreshTokenRepository        application_contracts.IRefreshTokenRepository
}

func NewLoginUserCommandHandler(unitOfWork application_contracts.IUnitOfWork,
	userRepository user_contracts.IUserRepository,
	encrypter application_contracts.IEncrypter,
	jwtProvider application_contracts.IJwtProvider,
	refreshTokenRepository application_contracts.IRefreshTokenRepository,
	externalCredentialsRepository application_contracts.IExternalCredentialsRepository,
) *LoginUserCommandHandler {
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
		userFromDb, err := handler.userRepository.GetUserByExternalId(ctx, command.CredentialsDto.ExternalId, handler.unitOfWork.GetQueryer())
		if err != nil {
			return customerrors.NewCustomError("LoginUserCommandHandler.GetUserByExternalId",
				"Problem with getting user with external id",
				err,
				customerrors.Unexpected)
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
			return customerrors.NewCustomError("LoginUserCommandHandler.Update",
				"Problem with updating user in the database",
				err,
				customerrors.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userFromDb.Id(), deviceId, string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return customerrors.NewCustomError("LoginUserCommandHandler.AssignNewToken",
				"Problem with assigning new token to user",
				newTokenErr,
				customerrors.Unexpected)
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
			command.CredentialsDto.ExternalId,
			command.CredentialsDto.MusicProvider,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		grantErr := handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
		if grantErr != nil {
			return customerrors.NewCustomError("LoginUserCommandHandler.Grant",
				"Problem with assigning external credentials to user",
				grantErr,
				customerrors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
