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

func (handler *CommandHandler) Handle(ctx context.Context, refreshToken *string) (*CommandResponse, error) {
	var response CommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.refreshTokenRepository.GetTokenByValue(ctx, *refreshToken, handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		if tokenFromDb.IsExpired() {
			return applicationErrors.NewApplicationError("expired token", nil, applicationErrors.Validation)
		}
		if !handler.encrypter.Verify(*refreshToken, []byte(tokenFromDb.RefreshToken())) {
			return applicationErrors.NewApplicationError("invalid refresh token", nil, applicationErrors.Unauthorized)
		}
		userFromDb, err := handler.userRepository.GetUserById(ctx, tokenFromDb.Id(), handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		newRefreshTokenValue := handler.encrypter.NewEncryptionKey()
		newRefreshToken := credentials.NewRefreshToken(userFromDb.Id(), tokenFromDb.DeviceId(), string(newRefreshTokenValue))
		newTokenErr := handler.refreshTokenRepository.UpdateToken(ctx, newRefreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return newTokenErr
		}
		response.RefreshToken = string(newRefreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return tokenErr
		}

		return nil
	})

	return &response, err
}
