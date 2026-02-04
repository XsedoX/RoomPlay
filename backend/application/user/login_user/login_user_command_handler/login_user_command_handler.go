package login_user_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/full_name"
)

type LoginUserCommandHandler struct {
	unitOfWork                    i_unit_of_work.IUnitOfWork
	userRepository                i_user_repository.IUserRepository
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
	encrypter                     i_encrypter.IEncrypter
	jwtProvider                   i_jwt_provider.IJwtProvider
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
}

func NewLoginUserCommandHandler(unitOfWork i_unit_of_work.IUnitOfWork,
	userRepository i_user_repository.IUserRepository,
	encrypter i_encrypter.IEncrypter,
	jwtProvider i_jwt_provider.IJwtProvider,
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository,
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository,
) *LoginUserCommandHandler {
	return &LoginUserCommandHandler{
		unitOfWork:                    unitOfWork,
		userRepository:                userRepository,
		encrypter:                     encrypter,
		jwtProvider:                   jwtProvider,
		internalCredentialsRepository: internalCredentialsRepository,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (handler *LoginUserCommandHandler) Handle(ctx context.Context, command *login_user_command.LoginUserCommand) (*login_user_command_response.LoginUserCommandResponse, error) {
	var response login_user_command_response.LoginUserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		userFromDb, err := handler.userRepository.GetUserByExternalId(ctx, command.CredentialsDto.ExternalId, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("LoginUserCommandHandler.GetUserByExternalId",
				"Problem with getting user with external id",
				err,
				custom_error_type.Unexpected)
		}

		userFromDb.ChangeFullName(full_name.NewFullName(command.Name, command.Surname))

		var deviceId device_id.DeviceId
		if command.DeviceDto.DeviceId == nil {
			deviceId = userFromDb.LoginWithNewDevice(command.DeviceDto.DeviceType)
		} else {
			userFromDb.ReloginWithKnownDevice(*command.DeviceDto.DeviceId)
			deviceId = *command.DeviceDto.DeviceId
		}
		response.DeviceId = deviceId
		err = handler.userRepository.Update(ctx, userFromDb, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("LoginUserCommandHandler.Update",
				"Problem with updating user in the database",
				err,
				custom_error_type.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		userSession := user_session.NewUserSession(userFromDb.Id(), deviceId)
		internalCredentials, internalCredsErr := internal_credentials.NewInternalCredentials(*userSession, string(refreshTokenValue))
		if internalCredsErr != nil {
			return internalCredsErr
		}

		newTokenErr := handler.internalCredentialsRepository.AssignNewToken(ctx, internalCredentials, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return custom_error.NewCustomError("LoginUserCommandHandler.AssignNewToken",
				"Problem with assigning new token to user",
				newTokenErr,
				custom_error_type.Unexpected)
		}
		response.RefreshToken = string(refreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return tokenErr
		}

		creds, externalCredsErr := external_credentials.NewExternalCredentials(userFromDb.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.ExternalId,
			command.CredentialsDto.MusicProvider,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		if externalCredsErr != nil {
			return externalCredsErr
		}

		grantErr := handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
		if grantErr != nil {
			return custom_error.NewCustomError("LoginUserCommandHandler.Grant",
				"Problem with assigning external credentials to user",
				grantErr,
				custom_error_type.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
