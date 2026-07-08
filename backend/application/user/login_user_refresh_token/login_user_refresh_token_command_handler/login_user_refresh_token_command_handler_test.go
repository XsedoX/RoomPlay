package login_user_refresh_token_command_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_jwt_provider"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_user_repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_internal_credentials_repository.MockInternalCredentialsRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*mock_user_repository.MockUserRepository,
	*mock_encrypter.MockEncrypter,
	*mock_jwt_provider.MockJwtProvider,
) {
	mockUserRepository := new(mock_user_repository.MockUserRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockInternalCredentialsRepository := new(mock_internal_credentials_repository.MockInternalCredentialsRepository)
	mockEncrypter := new(mock_encrypter.MockEncrypter)
	mockJwtProvider := new(mock_jwt_provider.MockJwtProvider)

	defer func() {
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockInternalCredentialsRepository.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
	}()

	return mockInternalCredentialsRepository, mockUoW, mockUserRepository, mockEncrypter, mockJwtProvider
}

func TestLoginUserRefreshTokenCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenGetTokenByValueFails", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		errToBeReturned := errors.New("could not get token")
		errCode := "LoginRefreshTokenCommandHandler.GetTokenByValue"
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(nil, errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *application_error.ApplicationError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
	})
	t.Run("ShouldReturnErrorWhenTokenIsExpired", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		errToBeReturned := application_error.NewApplicationError("LoginRefreshTokenCommandHandler.ExpiredToken",
			"Refresh token expired",
			nil,
			application_error_type.Unauthorized,
		)

		userSession := user_session.NewUserSession(user_id.NewUserId(), device_id.NewDeviceId())
		returnedRefreshToken := internal_credentials.HydrateInternalCredentials(
			*userSession,
			uuid.New().String(),
			time.Now().Add(-1*time.Hour).UTC(),
			time.Now().Add(-2*time.Hour).UTC())
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, errToBeReturned.Err)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unauthorized)
		assert.Equal(t, parsedErr.Title, errToBeReturned.Title)
		assert.Equal(t, parsedErr.Code, errToBeReturned.Code)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
	})
	t.Run("ShouldReturnErrorWhenGettingUserFromDbFails", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		userSession := user_session.NewUserSession(user_id.NewUserId(), device_id.NewDeviceId())
		returnedRefreshToken := internal_credentials.HydrateInternalCredentials(
			*userSession,
			uuid.New().String(),
			time.Now().Add(24*7*time.Hour).UTC(),
			time.Now().UTC())
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)
		userRepositoryErr := errors.New("userRepository error")
		errCodeToReturn := "LoginRefreshTokenCommandHandler.GetUserById"
		userId := userSession.UserId()
		mockUserRepository.On("GetUserById", mock.Anything, userId, mock.Anything).
			Return(nil, userRepositoryErr)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, userRepositoryErr)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unexpected)
		assert.Equal(t, parsedErr.Code, errCodeToReturn)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 2)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
	})
	t.Run("ShouldReturnErrorWhenAssigningNewTokenFails", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		userSession := user_session.NewUserSession(user_id.NewUserId(), device_id.NewDeviceId())
		returnedRefreshToken := internal_credentials.HydrateInternalCredentials(
			*userSession,
			uuid.New().String(),
			time.Now().Add(24*7*time.Hour).UTC(),
			time.Now().UTC())
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte(uuid.New().String()))
		devices := []device.Device{
			*device.HydrateDevice(returnedRefreshToken.DeviceId(),
				faker.Word(),
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC()),
		}
		userFromDb := user.HydrateUser(userSession.UserId(),
			faker.FirstName(),
			faker.LastName(),
			nil,
			nil,
			devices,
			nil,
		)
		userId := userSession.UserId()
		mockUserRepository.On("GetUserById", mock.Anything, userId, mock.Anything).
			Return(userFromDb, nil)
		assignTokenErr := errors.New("assignToken error")
		errCodeToReturn := "LoginRefreshTokenCommandHandler.AssignNewToken"
		mockInternalCredentialsRepository.
			On("AssignNewToken", mock.Anything, mock.AnythingOfType("*internal_credentials.InternalCredentials"), mock.Anything).
			Return(assignTokenErr)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, assignTokenErr)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unexpected)
		assert.Equal(t, parsedErr.Code, errCodeToReturn)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 3)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
	})
	t.Run("ShouldReturnErrorWhenGenerateNewTokenFails", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		userSession := user_session.NewUserSession(user_id.NewUserId(), device_id.NewDeviceId())
		returnedRefreshToken := internal_credentials.HydrateInternalCredentials(
			*userSession,
			uuid.New().String(),
			time.Now().Add(24*7*time.Hour).UTC(),
			time.Now().UTC())
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte(uuid.New().String()))
		devices := []device.Device{
			*device.HydrateDevice(returnedRefreshToken.DeviceId(),
				faker.Word(),
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC()),
		}
		userFromDb := user.HydrateUser(userSession.UserId(),
			faker.FirstName(),
			faker.LastName(),
			nil,
			nil,
			devices,
			nil,
		)
		userSessionFromToken := returnedRefreshToken.UserSession()
		mockUserRepository.On("GetUserById", mock.Anything, userSessionFromToken.UserId(), mock.Anything).
			Return(userFromDb, nil)
		mockInternalCredentialsRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*internal_credentials.InternalCredentials"), mock.Anything).Return(nil)
		generateTokenErr := errors.New("generateToken error")
		errCodeToReturn := "LoginRefreshTokenCommandHandler.GenerateToken"
		mockJwtProvider.On("GenerateToken", userSessionFromToken.UserId()).Return("", generateTokenErr)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.Empty(t, resp)
		assert.Error(t, handlerErr)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(handlerErr, &parsedErr))
		assert.Equal(t, parsedErr.Err, generateTokenErr)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unexpected)
		assert.Equal(t, parsedErr.Code, errCodeToReturn)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 3)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockJwtProvider.AssertNumberOfCalls(t, "GenerateToken", 1)
	})
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockInternalCredentialsRepository,
			mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider := setupMocks(t)

		handler := NewLoginUserRefreshTokenCommandHandler(mockInternalCredentialsRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockUserRepository,
		)
		tokenCommand := uuid.New().String()
		mockUoW.On("GetQueryer").Return(nil)
		userSession := user_session.NewUserSession(user_id.NewUserId(), device_id.NewDeviceId())
		returnedRefreshToken := internal_credentials.HydrateInternalCredentials(
			*userSession,
			uuid.New().String(),
			time.Now().Add(24*7*time.Hour).UTC(),
			time.Now().UTC())
		mockInternalCredentialsRepository.
			On("GetTokenByValue", mock.Anything, tokenCommand, mock.Anything).
			Return(returnedRefreshToken, nil)
		refreshTokenToReturn := uuid.New().String()
		mockEncrypter.On("NewEncryptionKey").Return([]byte(refreshTokenToReturn))
		devices := []device.Device{
			*device.HydrateDevice(returnedRefreshToken.DeviceId(),
				faker.Word(),
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC()),
		}
		userFromDb := user.HydrateUser(userSession.UserId(),
			faker.FirstName(),
			faker.LastName(),
			nil,
			nil,
			devices,
			nil,
		)
		userSessionFromToken := returnedRefreshToken.UserSession()
		mockUserRepository.On("GetUserById", mock.Anything, userSessionFromToken.UserId(), mock.Anything).
			Return(userFromDb, nil)
		var passedRefreshToken *internal_credentials.InternalCredentials
		mockInternalCredentialsRepository.
			On("AssignNewToken", mock.Anything, mock.AnythingOfType("*internal_credentials.InternalCredentials"), mock.Anything).
			Run(func(args mock.Arguments) {
				passedRefreshToken = args.Get(1).(*internal_credentials.InternalCredentials)
			}).
			Return(nil)
		accessTokenToReturn := uuid.New().String()
		mockJwtProvider.On("GenerateToken", userSessionFromToken.UserId()).Return(accessTokenToReturn, nil)

		resp, handlerErr := handler.Handle(context.Background(), &tokenCommand)

		assert.NoError(t, handlerErr)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.AccessToken, accessTokenToReturn)
		assert.Equal(t, resp.RefreshToken, refreshTokenToReturn)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 3)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "GetTokenByValue", 1)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockJwtProvider.AssertNumberOfCalls(t, "GenerateToken", 1)
		mockEncrypter.AssertNumberOfCalls(t, "NewEncryptionKey", 1)
		mockInternalCredentialsRepository.AssertNumberOfCalls(t, "AssignNewToken", 1)
		passedUserSession := passedRefreshToken.UserSession()
		assert.Equal(t, passedUserSession.UserId(), userFromDb.Id())
		assert.Equal(t, passedUserSession.DeviceId(), devices[0].Id())
		assert.Equal(t, passedRefreshToken.RefreshToken(), refreshTokenToReturn)
	})
}
