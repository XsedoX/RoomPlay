package login_user_refresh_token_command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/tests/infrustructure/authentication"
	persistance2 "xsedox.com/main/tests/infrustructure/persistance"
)

func TestLoginUserRefreshTokenCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenGetTokenByValueFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		handler := NewLoginUserRefreshTokenCommandHandler(mockRefreshTokenRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		errToBeReturned := errors.New("could not get token")
		errCode := "LoginRefreshTokenCommandHandler.GetTokenByValue"
		mockRefreshTokenRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(nil, errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
		mockUoW.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
	})
	t.Run("ShouldReturnErrorWhenTokenIsExpired", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		handler := NewLoginUserRefreshTokenCommandHandler(mockRefreshTokenRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		errToBeReturned := custom_errors.NewCustomError("LoginRefreshTokenCommandHandler.ExpiredToken",
			"Refresh token expired",
			nil,
			custom_errors.Unauthorized)
		returnedRefreshToken := credentials.HydrateRefreshToken(
			user.Id(uuid.New()),
			user.DeviceId(uuid.New()),
			uuid.New().String(),
			time.Now().Add(-1*time.Hour).UTC(),
			time.Now().Add(-2*time.Hour).UTC())
		mockRefreshTokenRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, errToBeReturned.Err)
		assert.Equal(t, parsedErr.ErrorType, custom_errors.Unauthorized)
		assert.Equal(t, parsedErr.Title, errToBeReturned.Title)
		assert.Equal(t, parsedErr.Code, errToBeReturned.Code)
		mockUoW.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
	})
	t.Run("ShouldReturnErrorWhenGettingUserFromDbFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		handler := NewLoginUserRefreshTokenCommandHandler(mockRefreshTokenRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		returnedRefreshToken := credentials.HydrateRefreshToken(
			user.Id(uuid.New()),
			user.DeviceId(uuid.New()),
			uuid.New().String(),
			time.Now().Add(24*7*time.Hour).UTC(),
			time.Now().UTC())
		mockRefreshTokenRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)
		userRepositoryErr := errors.New("userRepository error")
		errCodeToReturn := "LoginRefreshTokenCommandHandler.GetUserById"
		mockUserRepository.On("GetUserById", mock.Anything, returnedRefreshToken.Id(), mock.Anything).
			Return(nil, userRepositoryErr)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, userRepositoryErr)
		assert.Equal(t, parsedErr.ErrorType, custom_errors.Unexpected)
		assert.Equal(t, parsedErr.Code, errCodeToReturn)
		mockUoW.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 2)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
	})
}
