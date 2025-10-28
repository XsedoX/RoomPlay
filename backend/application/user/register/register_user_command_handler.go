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
	userRepository         contracts.IUserRepository
	credentialsRepository  contracts.ICredentialsRepository
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
	jwtProvider            contracts.IJwtProvider
	encrypter              contracts.IEncrypter
}

func NewRegisterUserCommandHandler(userRepository contracts.IUserRepository,
	unitOfWork contracts.IUnitOfWork,
	credsRepository contracts.ICredentialsRepository,
	jwtProvider contracts.IJwtProvider,
	refreshTokenRepository contracts.IRefreshTokenRepository,
	encrypter contracts.IEncrypter) *UserCommandHandler {
	return &UserCommandHandler{
		userRepository:         userRepository,
		credentialsRepository:  credsRepository,
		unitOfWork:             unitOfWork,
		jwtProvider:            jwtProvider,
		refreshTokenRepository: refreshTokenRepository,
		encrypter:              encrypter,
	}
}
func (handler *UserCommandHandler) Handle(ctx context.Context, command *UserCommand) (*UserCommandResponse, *applicationErrors.ApplicationError) {
	var response UserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		deviceEnt := device.NewDevice(command.DeviceType)
		userAgg := user.FirstLogin(command.ExternalId, command.Name, command.Surname, *deviceEnt)
		creds := credentials.NewExternalCredentials(userAgg.Id(),
			command.CredentialsDto.AccessToken,
			command.CredentialsDto.RefreshToken,
			command.CredentialsDto.Scopes,
			command.CredentialsDto.AccessTokenExpiresAt,
			command.CredentialsDto.RefreshTokenExpiresAt)

		refreshTokenValue := handler.encrypter.NewEncryptionKey()
		refreshToken := credentials.NewRefreshToken(userAgg.Id(), deviceEnt.Id(), refreshTokenValue)

		err := handler.userRepository.AddUser(ctx, userAgg, handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userAgg.Id())
		if tokenErr != nil {
			return tokenErr
		}

		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return newTokenErr
		}
		response.RefreshToken = string(refreshTokenValue)
		response.DeviceId = deviceEnt.Id()
		return handler.credentialsRepository.Grant(ctx, creds, handler.unitOfWork.GetQueryer())
	})
	respErr := applicationErrors.NewApplicationError("problem with executing transaction", err, applicationErrors.ErrInfrastructure)
	return &response, respErr
}
