package login_with_refresh_token

import (
	"context"
	"errors"
	"time"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/services"
)

type CommandHandler struct {
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
	encrypter              contracts.IEncrypter
	jwtProvider            contracts.IJwtProvider
}

func NewRefreshTokenCommandHandler(refreshToken contracts.IRefreshTokenRepository,
	unitOfWork contracts.IUnitOfWork,
	encrypter contracts.IEncrypter,
	jwtProvider contracts.IJwtProvider) *CommandHandler {
	return &CommandHandler{
		refreshTokenRepository: refreshToken,
		unitOfWork:             unitOfWork,
		encrypter:              encrypter,
		jwtProvider:            jwtProvider}
}

func (handler *CommandHandler) Handle(ctx context.Context, refreshToken string) (*services.CommandResponse, error) {
	var response services.CommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.refreshTokenRepository.GetTokenByValue(ctx, refreshToken, handler.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		if tokenFromDb.ExpirationTime().Sub(time.Now().UTC()) <= 0 || !handler.encrypter.Verify(refreshToken, tokenFromDb.RefreshToken()) {
			return errors.New("refresh token expired")
		}
		return nil
	})
	return &response, err
}
