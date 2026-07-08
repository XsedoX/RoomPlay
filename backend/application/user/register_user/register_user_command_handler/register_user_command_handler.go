package register_user_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type RegisterUserCommandHandler struct {
	userRepository                i_user_repository.IUserRepository
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
	unitOfWork                    i_unit_of_work.IUnitOfWork
	jwtProvider                   i_jwt_provider.IJwtProvider
	encrypter                     i_encrypter.IEncrypter
}

func NewRegisterUserCommandHandler(userRepository i_user_repository.IUserRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	credsRepository i_external_credentials_repository.IExternalCredentialsRepository,
	jwtProvider i_jwt_provider.IJwtProvider,
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository,
	encrypter i_encrypter.IEncrypter,
) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{
		userRepository:                userRepository,
		externalCredentialsRepository: credsRepository,
		unitOfWork:                    unitOfWork,
		jwtProvider:                   jwtProvider,
		internalCredentialsRepository: internalCredentialsRepository,
		encrypter:                     encrypter,
	}
}

func (handler *RegisterUserCommandHandler) Handle(ctx context.Context, command *register_user_command.RegisterUserCommand) (*register_user_command_response.RegisterUserCommandResponse, error) {
	var response register_user_command_response.RegisterUserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		userAgg := user.NewUser(command.Name, command.Surname, command.DeviceType)
		deviceEnt := userAgg.GetMostRecentDevice()
		err := handler.userRepository.Add(ctx, userAgg, handler.unitOfWork.GetQueryer())
		if err != nil {
			return application_error.NewApplicationError("RegisterUserCommandHandler.UserRepository.Add",
				"Adding user problem",
				err,
				application_error_type.Unexpected)
		}
		response.DeviceId = deviceEnt.Id()

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userAgg.Id())
		if tokenErr != nil {
			return application_error.NewApplicationError("RegisterUserCommandHandler.GenerateToken",
				"Access token generation problem",
				tokenErr,
				application_error_type.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		userSession := user_session.NewUserSession(userAgg.Id(), deviceEnt.Id())
		internalCredentials, internalCredentialsErr := internal_credentials.NewInternalCredentials(*userSession, string(refreshTokenValue))
		if internalCredentialsErr != nil {
			return internalCredentialsErr
		}
		newTokenErr := handler.internalCredentialsRepository.AssignNewToken(ctx, internalCredentials, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return application_error.NewApplicationError("RegisterUserCommandHandler.AssignNewToken",
				"Access token generation problem",
				newTokenErr,
				application_error_type.Unexpected)
		}
		response.RefreshToken = string(refreshTokenValue)

		creds, externalCredentialsErr := external_credentials.NewExternalCredentials(userAgg.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.ExternalId,
			command.CredentialsDto.MusicProvider,
			command.CredentialsDto.AccessTokenExpiresAtUtc,
			command.CredentialsDto.RefreshTokenExpiresAtUtc)
		if externalCredentialsErr != nil {
			return externalCredentialsErr
		}
		grantErr := handler.externalCredentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
		if grantErr != nil {
			return application_error.NewApplicationError("RegisterUserCommandHandler.Grant",
				"Problem with granting external credentials.",
				grantErr,
				application_error_type.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, err
}
