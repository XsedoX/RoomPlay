package login_refresh_token

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
)

type CommandHandler struct {
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
	encrypter              contracts.IEncrypter
	userRepository         contracts.IUserRepository
	jwtProvider            contracts.IJwtProvider
}

func NewLoginRefreshTokenCommandHandler(refreshToken contracts.IRefreshTokenRepository,
	unitOfWork contracts.IUnitOfWork,
	encrypter contracts.IEncrypter,
	jwtProvider contracts.IJwtProvider,
	userRepository contracts.IUserRepository) *CommandHandler {
	return &CommandHandler{
		refreshTokenRepository: refreshToken,
		unitOfWork:             unitOfWork,
		encrypter:              encrypter,
		userRepository:         userRepository,
		jwtProvider:            jwtProvider}
}

func (handler *CommandHandler) Handle(ctx context.Context, command *string) (*CommandResponse, error) {
	var response CommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.refreshTokenRepository.GetTokenByValue(ctx, *command, handler.unitOfWork.GetQueryer())
		if err != nil {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.GetTokenByValue",
				"Couldn't fetch token from the database.",
				err,
				applicationErrors.Unexpected)
		}
		if tokenFromDb.IsExpired() {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.ExpiredToken",
				"Refresh token expired",
				nil,
				applicationErrors.Unauthorized)
		}
		if !handler.encrypter.Verify(*command, []byte(tokenFromDb.RefreshToken())) {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.InvalidRefreshToken",
				"The token provided is invalid",
				nil,
				applicationErrors.Unauthorized)
		}
		userFromDb, err := handler.userRepository.GetUserById(ctx, tokenFromDb.Id(), handler.unitOfWork.GetQueryer())
		if err != nil {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.GetUserById",
				"Couldn't fetch user from the database.",
				err,
				applicationErrors.Unexpected)
		}
		newRefreshTokenValue := handler.encrypter.NewEncryptionKey()
		newRefreshToken := credentials.NewRefreshToken(userFromDb.Id(), tokenFromDb.DeviceId(), string(newRefreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, newRefreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.AssignNewToken",
				"Couldn't assign a new refresh token to a user.",
				newTokenErr,
				applicationErrors.Unexpected)
		}
		response.RefreshToken = string(newRefreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return applicationErrors.NewApplicationError("LoginRefreshTokenCommandHandler.GenerateToken",
				"Couldn't generate a new access token for a user.",
				newTokenErr,
				applicationErrors.Unexpected)
		}

		return nil
	})

	return &response, err
}
