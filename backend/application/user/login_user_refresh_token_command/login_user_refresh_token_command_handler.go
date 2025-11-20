package login_user_refresh_token_command

import (
	"context"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts2 "xsedox.com/main/application/user/contracts"
	"xsedox.com/main/domain/credentials"
)

type LoginUserRefreshTokenCommandHandler struct {
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
	encrypter              contracts.IEncrypter
	userRepository         contracts2.IUserRepository
	jwtProvider            contracts.IJwtProvider
}

func NewLoginUserRefreshTokenCommandHandler(refreshToken contracts.IRefreshTokenRepository,
	unitOfWork contracts.IUnitOfWork,
	encrypter contracts.IEncrypter,
	jwtProvider contracts.IJwtProvider,
	userRepository contracts2.IUserRepository) *LoginUserRefreshTokenCommandHandler {
	return &LoginUserRefreshTokenCommandHandler{
		refreshTokenRepository: refreshToken,
		unitOfWork:             unitOfWork,
		encrypter:              encrypter,
		userRepository:         userRepository,
		jwtProvider:            jwtProvider}
}

func (handler *LoginUserRefreshTokenCommandHandler) Handle(ctx context.Context, command *string) (*LoginUserRefreshTokenCommandResponse, error) {
	var response LoginUserRefreshTokenCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.refreshTokenRepository.GetTokenByValue(ctx, *command, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.GetTokenByValue",
				"Couldn't fetch token from the database.",
				err,
				custom_errors.Unauthorized)
		}
		if tokenFromDb.IsExpired() {
			return custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.ExpiredToken",
				"Refresh token expired",
				nil,
				custom_errors.Unauthorized)
		}
		userFromDb, err := handler.userRepository.GetUserById(ctx, tokenFromDb.Id(), handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.GetUserById",
				"Couldn't fetch user from the database.",
				err,
				custom_errors.Unexpected)
		}
		newRefreshTokenValue := handler.encrypter.NewEncryptionKey()
		newRefreshToken := credentials.NewRefreshToken(userFromDb.Id(), tokenFromDb.DeviceId(), string(newRefreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, newRefreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.AssignNewToken",
				"Couldn't assign a new refresh token to a user.",
				newTokenErr,
				custom_errors.Unexpected)
		}
		response.RefreshToken = string(newRefreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.GenerateToken",
				"Couldn't generate a new access token for a user.",
				newTokenErr,
				custom_errors.Unexpected)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &response, err
}
