package register_user

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/test_helpers/infrustructure_test/authentication_mocks"
	"xsedox.com/main/test_helpers/infrustructure_test/persistance_mocks"
)

func TestRegisterUserCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenAddingUserFails", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		mockJwtProvider := new(authentication_mocks.MockJwtProvider)
		mockExternalCredentialsRepository := new(persistance_mocks.MockExternalCredentialsRepository)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockEncrypter,
		)

		command := RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			ExternalId: uuid.New().String(),
			DeviceType: user.Mobile,
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "scope",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAtUtc: time.Now().Add(24 * time.Hour),
				IssuedAt:                 time.Now(),
			},
		}

		errToBeReturned := errors.New("could not add user")
		errCode := "RegisterUserCommandHandler.UserRepository.Add"

		mockUoW.On("GetQueryer").Return(nil)

		mockUserRepository.On("Add", mock.Anything, mock.AnythingOfType("*user.User"), mock.Anything).
			Return(errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)

		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Add", 1)
	})

	t.Run("ShouldReturnErrorWhenGenerateTokenFails", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		mockJwtProvider := new(authentication_mocks.MockJwtProvider)
		mockExternalCredentialsRepository := new(persistance_mocks.MockExternalCredentialsRepository)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockEncrypter,
		)

		command := RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			ExternalId: uuid.New().String(),
			DeviceType: user.Mobile,
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "scope",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAtUtc: time.Now().Add(24 * time.Hour),
				IssuedAt:                 time.Now(),
			},
		}

		errToBeReturned := errors.New("token generation failed")
		errCode := "RegisterUserCommandHandler.GenerateToken"

		mockUoW.On("GetQueryer").Return(nil)

		mockUserRepository.On("Add", mock.Anything, mock.AnythingOfType("*user.User"), mock.Anything).
			Return(nil)

		mockJwtProvider.On("GenerateToken", mock.Anything).
			Return("", errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)

		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
	})

	t.Run("ShouldReturnErrorWhenAssigningNewTokenFails", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		mockJwtProvider := new(authentication_mocks.MockJwtProvider)
		mockExternalCredentialsRepository := new(persistance_mocks.MockExternalCredentialsRepository)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockEncrypter,
		)

		command := RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			ExternalId: uuid.New().String(),
			DeviceType: user.Mobile,
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "scope",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAtUtc: time.Now().Add(24 * time.Hour),
				IssuedAt:                 time.Now(),
			},
		}

		errToBeReturned := errors.New("assign token failed")
		errCode := "RegisterUserCommandHandler.AssignNewToken"

		mockUoW.On("GetQueryer").Return(nil)

		mockUserRepository.On("Add", mock.Anything, mock.AnythingOfType("*user.User"), mock.Anything).
			Return(nil)

		mockJwtProvider.On("GenerateToken", mock.Anything).
			Return("access_token", nil)

		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh_token_key"))

		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Return(errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	})

	t.Run("ShouldReturnErrorWhenGrantingExternalCredentialsFails", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		mockJwtProvider := new(authentication_mocks.MockJwtProvider)
		mockExternalCredentialsRepository := new(persistance_mocks.MockExternalCredentialsRepository)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockEncrypter,
		)

		command := RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			ExternalId: uuid.New().String(),
			DeviceType: user.Mobile,
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "scope",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAtUtc: time.Now().Add(24 * time.Hour),
				IssuedAt:                 time.Now(),
			},
		}

		errToBeReturned := errors.New("grant credentials failed")
		errCode := "RegisterUserCommandHandler.Grant"

		mockUoW.On("GetQueryer").Return(nil)

		mockUserRepository.On("Add", mock.Anything, mock.AnythingOfType("*user.User"), mock.Anything).
			Return(nil)

		mockJwtProvider.On("GenerateToken", mock.Anything).
			Return("access_token", nil)

		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh_token_key"))

		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Return(nil)

		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.External"), mock.Anything).
			Return(errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	})

	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		mockJwtProvider := new(authentication_mocks.MockJwtProvider)
		mockExternalCredentialsRepository := new(persistance_mocks.MockExternalCredentialsRepository)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockEncrypter,
		)

		command := RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			ExternalId: uuid.New().String(),
			DeviceType: user.Mobile,
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "scope",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAtUtc: time.Now().Add(24 * time.Hour),
				IssuedAt:                 time.Now(),
			},
		}

		mockUoW.On("GetQueryer").Return(nil)

		var capturedUser *user.User
		mockUserRepository.On("Add", mock.Anything, mock.AnythingOfType("*user.User"), mock.Anything).
			Run(func(args mock.Arguments) {
				capturedUser = args.Get(1).(*user.User)
			}).
			Return(nil)

		expectedAccessToken := "generated_access_token"
		mockJwtProvider.On("GenerateToken", mock.Anything).
			Return(expectedAccessToken, nil)

		expectedRefreshTokenKey := "generated_refresh_token_key"
		mockEncrypter.On("NewEncryptionKey").Return([]byte(expectedRefreshTokenKey))

		var capturedRefreshToken *credentials.RefreshToken
		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Run(func(args mock.Arguments) {
				capturedRefreshToken = args.Get(1).(*credentials.RefreshToken)
			}).
			Return(nil)

		var capturedCreds *credentials.External
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.External"), mock.Anything).
			Run(func(args mock.Arguments) {
				capturedCreds = args.Get(1).(*credentials.External)
			}).
			Return(nil)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.NoError(t, handlerErr)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedAccessToken, resp.AccessToken)
		assert.Equal(t, expectedRefreshTokenKey, resp.RefreshToken)
		// We can't easily assert DeviceId without capturing the user first, as it's generated inside.
		// But we can check if it matches the captured user's device.

		assert.NotNil(t, capturedUser)
		assert.Equal(t, command.Name, capturedUser.FullName().Name())
		assert.Equal(t, command.Surname, capturedUser.FullName().Surname())
		assert.Equal(t, command.ExternalId, capturedUser.ExternalId())
		assert.Equal(t, command.DeviceType, capturedUser.GetMostRecentDevice().DeviceType())
		assert.Equal(t, capturedUser.GetMostRecentDevice().Id(), resp.DeviceId)

		assert.NotNil(t, capturedRefreshToken)
		assert.Equal(t, capturedUser.Id(), capturedRefreshToken.Id())
		assert.Equal(t, capturedUser.GetMostRecentDevice().Id(), capturedRefreshToken.DeviceId())
		assert.Equal(t, expectedRefreshTokenKey, capturedRefreshToken.RefreshToken())

		assert.NotNil(t, capturedCreds)
		assert.Equal(t, capturedUser.Id(), capturedCreds.Id())
		assert.Equal(t, command.CredentialsDto.AccessToken, capturedCreds.AccessToken())
		assert.Equal(t, command.CredentialsDto.RefreshToken, capturedCreds.RefreshToken())
		assert.Equal(t, strings.Split(command.CredentialsDto.Scopes, " "), capturedCreds.Scopes())
		assert.Equal(t, command.CredentialsDto.AccessTokenExpiresAtUtc, capturedCreds.AccessTokenExpiresAtUtc())
		assert.Equal(t, command.CredentialsDto.RefreshTokenExpiresAtUtc, capturedCreds.RefreshTokenExpiresAtUtc())

		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	})
}
