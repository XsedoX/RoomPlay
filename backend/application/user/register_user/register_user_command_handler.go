package register_user

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	contracts2 "github.com/XsedoX/RoomPlay/application/user/contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
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
	encrypter contracts.IEncrypter,
) *RegisterUserCommandHandler {
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
			return customerrors.NewCustomError("RegisterUserCommandHandler.UserRepository.Add",
				"Adding user problem",
				err,
				customerrors.Unexpected)
		}
		response.DeviceId = deviceEnt.Id()

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userAgg.Id())
		if tokenErr != nil {
			return customerrors.NewCustomError("RegisterUserCommandHandler.GenerateToken",
				"Access token generation problem",
				tokenErr,
				customerrors.Unexpected)
		}

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userAgg.Id(), deviceEnt.Id(), string(refreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return customerrors.NewCustomError("RegisterUserCommandHandler.AssignNewToken",
				"Access token generation problem",
				newTokenErr,
				customerrors.Unexpected)
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
			return customerrors.NewCustomError("RegisterUserCommandHandler.Grant",
				"Problem with granting external credentials.",
				grantErr,
				customerrors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, err
}
