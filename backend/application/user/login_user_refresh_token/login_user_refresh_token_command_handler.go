package login_user_refresh_token

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
)

type LoginUserRefreshTokenCommandHandler struct {
	refreshTokenRepository application_contracts.IRefreshTokenRepository
	unitOfWork             application_contracts.IUnitOfWork
	encrypter              application_contracts.IEncrypter
	userRepository         user_contracts.IUserRepository
	jwtProvider            application_contracts.IJwtProvider
}

func NewLoginUserRefreshTokenCommandHandler(refreshToken application_contracts.IRefreshTokenRepository,
	unitOfWork application_contracts.IUnitOfWork,
	encrypter application_contracts.IEncrypter,
	jwtProvider application_contracts.IJwtProvider,
	userRepository user_contracts.IUserRepository,
) *LoginUserRefreshTokenCommandHandler {
	return &LoginUserRefreshTokenCommandHandler{
		refreshTokenRepository: refreshToken,
		unitOfWork:             unitOfWork,
		encrypter:              encrypter,
		userRepository:         userRepository,
		jwtProvider:            jwtProvider,
	}
}

func (handler *LoginUserRefreshTokenCommandHandler) Handle(ctx context.Context, command *string) (*LoginUserRefreshTokenCommandResponse, error) {
	var response LoginUserRefreshTokenCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.refreshTokenRepository.GetTokenByValue(ctx, *command, handler.unitOfWork.GetQueryer())
		if err != nil {
			return customerrors.NewCustomError("LoginRefreshTokenCommandHandler.GetTokenByValue",
				"Couldn't fetch token from the database.",
				err,
				customerrors.Unauthorized)
		}
		if tokenFromDb.IsExpired() {
			return customerrors.NewCustomError("LoginRefreshTokenCommandHandler.ExpiredToken",
				"Refresh token expired",
				nil,
				customerrors.Unauthorized)
		}
		userFromDb, err := handler.userRepository.GetUserById(ctx, tokenFromDb.Id(), handler.unitOfWork.GetQueryer())
		if err != nil {
			return customerrors.NewCustomError("LoginRefreshTokenCommandHandler.GetUserById",
				"Couldn't fetch user from the database.",
				err,
				customerrors.Unexpected)
		}
		newRefreshTokenValue := handler.encrypter.NewEncryptionKey()
		newRefreshToken := credentials.NewRefreshToken(userFromDb.Id(), tokenFromDb.DeviceId(), string(newRefreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.AssignNewToken(ctx, newRefreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return customerrors.NewCustomError("LoginRefreshTokenCommandHandler.AssignNewToken",
				"Couldn't assign a new refresh token to a user.",
				newTokenErr,
				customerrors.Unexpected)
		}
		response.RefreshToken = string(newRefreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return customerrors.NewCustomError("LoginRefreshTokenCommandHandler.GenerateToken",
				"Couldn't generate a new access token for a user.",
				tokenErr,
				customerrors.Unexpected)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &response, err
}
