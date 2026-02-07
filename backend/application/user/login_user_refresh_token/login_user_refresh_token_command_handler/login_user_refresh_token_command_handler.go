package login_user_refresh_token_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token/login_user_refresh_token_command_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
)

type LoginUserRefreshTokenCommandHandler struct {
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
	unitOfWork                    i_unit_of_work.IUnitOfWork
	encrypter                     i_encrypter.IEncrypter
	userRepository                i_user_repository.IUserRepository
	jwtProvider                   i_jwt_provider.IJwtProvider
}

func NewLoginUserRefreshTokenCommandHandler(internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	encrypter i_encrypter.IEncrypter,
	jwtProvider i_jwt_provider.IJwtProvider,
	userRepository i_user_repository.IUserRepository,
) *LoginUserRefreshTokenCommandHandler {
	return &LoginUserRefreshTokenCommandHandler{
		internalCredentialsRepository: internalCredentialsRepository,
		unitOfWork:                    unitOfWork,
		encrypter:                     encrypter,
		userRepository:                userRepository,
		jwtProvider:                   jwtProvider,
	}
}

func (handler *LoginUserRefreshTokenCommandHandler) Handle(ctx context.Context, command *string) (*login_user_refresh_token_command_response.LoginUserRefreshTokenCommandResponse, error) {
	var response login_user_refresh_token_command_response.LoginUserRefreshTokenCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		tokenFromDb, err := handler.internalCredentialsRepository.GetTokenByValue(ctx, *command, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("LoginRefreshTokenCommandHandler.GetTokenByValue",
				"Couldn't fetch token from the database.",
				err,
				custom_error_type.Unauthorized,
			)
		}
		if tokenFromDb.IsExpired() {
			return custom_error.NewCustomError("LoginRefreshTokenCommandHandler.ExpiredToken",
				"Refresh token expired",
				nil,
				custom_error_type.Unauthorized,
			)
		}
		userFromDb, err := handler.userRepository.GetUserById(ctx, tokenFromDb.UserId(), handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("LoginRefreshTokenCommandHandler.GetUserById",
				"Couldn't fetch user from the database.",
				err,
				custom_error_type.Unexpected,
			)
		}
		newRefreshTokenValue := handler.encrypter.NewEncryptionKey()
		newRefreshToken, refreshTokenErr := internal_credentials.NewInternalCredentials(tokenFromDb.UserSession(), string(newRefreshTokenValue))
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		newTokenErr := handler.internalCredentialsRepository.AssignNewToken(ctx, newRefreshToken, handler.unitOfWork.GetQueryer())
		if newTokenErr != nil {
			return custom_error.NewCustomError("LoginRefreshTokenCommandHandler.AssignNewToken",
				"Couldn't assign a new refresh token to a user.",
				newTokenErr,
				custom_error_type.Unexpected,
			)
		}
		response.RefreshToken = string(newRefreshTokenValue)

		var tokenErr error
		response.AccessToken, tokenErr = handler.jwtProvider.GenerateToken(userFromDb.Id())
		if tokenErr != nil {
			return custom_error.NewCustomError("LoginRefreshTokenCommandHandler.GenerateToken",
				"Couldn't generate a new access token for a user.",
				tokenErr,
				custom_error_type.Unexpected,
			)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &response, err
}
