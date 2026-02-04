package register_user_command_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_jwt_provider"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_user_repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_user_repository.MockUserRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*mock_internal_credentials_repository.MockInternalCredentialsRepository,
	*mock_jwt_provider.MockJwtProvider,
	*mock_encrypter.MockEncrypter,
	*mock_external_credentials_repository.MockExternalCredentialsRepository,
) {
	mockUserRepository := new(mock_user_repository.MockUserRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockInternalCredentialsRepository := new(mock_internal_credentials_repository.MockInternalCredentialsRepository)
	mockEncrypter := new(mock_encrypter.MockEncrypter)
	mockJwtProvider := new(mock_jwt_provider.MockJwtProvider)
	mockExternalCredentialsRepository := new(mock_external_credentials_repository.MockExternalCredentialsRepository)

	defer func() {
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockInternalCredentialsRepository.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	}()

	return mockUserRepository, mockUoW, mockInternalCredentialsRepository, mockJwtProvider, mockEncrypter, mockExternalCredentialsRepository
}

func TestRegisterUserCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenAddingUserFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockInternalCredentialsRepository,
			mockJwtProvider,
			mockEncrypter,
			mockExternalCredentialsRepository := setupMocks(t)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockInternalCredentialsRepository,
			mockEncrypter,
		)

		command := register_user_command.RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			DeviceType: device_type.Mobile,
			CredentialsDto: register_user_command.CredentialsDto{
				ExternalId:               uuid.New().String(),
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
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
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)

		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Add", 1)
	})

	t.Run("ShouldReturnErrorWhenGenerateTokenFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockInternalCredentialsRepository,
			mockJwtProvider,
			mockEncrypter,
			mockExternalCredentialsRepository := setupMocks(t)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockInternalCredentialsRepository,
			mockEncrypter,
		)

		command := register_user_command.RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			DeviceType: device_type.Mobile,
			CredentialsDto: register_user_command.CredentialsDto{
				ExternalId:               uuid.New().String(),
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
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
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
	})

	t.Run("ShouldReturnErrorWhenAssigningNewTokenFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockInternalCredentialsRepository,
			mockJwtProvider,
			mockEncrypter,
			mockExternalCredentialsRepository := setupMocks(t)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockInternalCredentialsRepository,
			mockEncrypter,
		)

		command := register_user_command.RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			DeviceType: device_type.Mobile,
			CredentialsDto: register_user_command.CredentialsDto{
				ExternalId:               uuid.New().String(),
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
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

		mockInternalCredentialsRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Return(errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
	})

	t.Run("ShouldReturnErrorWhenGrantingExternalCredentialsFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockInternalCredentialsRepository,
			mockJwtProvider,
			mockEncrypter,
			mockExternalCredentialsRepository := setupMocks(t)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockInternalCredentialsRepository,
			mockEncrypter,
		)

		command := register_user_command.RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			DeviceType: device_type.Mobile,
			CredentialsDto: register_user_command.CredentialsDto{
				ExternalId:               uuid.New().String(),
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
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

		mockInternalCredentialsRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Return(nil)

		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.ExternalCredentials"), mock.Anything).
			Return(errToBeReturned)

		resp, handlerErr := handler.Handle(context.Background(), &command)

		assert.Nil(t, resp)
		assert.Error(t, handlerErr)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(handlerErr, &customErr))
		assert.Equal(t, errCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, errToBeReturned)
	})

	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockInternalCredentialsRepository,
			mockJwtProvider,
			mockEncrypter,
			mockExternalCredentialsRepository := setupMocks(t)

		handler := NewRegisterUserCommandHandler(
			mockUserRepository,
			mockUoW,
			mockExternalCredentialsRepository,
			mockJwtProvider,
			mockInternalCredentialsRepository,
			mockEncrypter,
		)

		command := register_user_command.RegisterUserCommand{
			Name:       faker.FirstName(),
			Surname:    faker.LastName(),
			DeviceType: device_type.Mobile,
			CredentialsDto: register_user_command.CredentialsDto{
				ExternalId:               uuid.New().String(),
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
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

		var capturedRefreshToken *internal_credentials.InternalCredentials
		mockInternalCredentialsRepository.On("AssignNewToken", mock.Anything, mock.AnythingOfType("*credentials.RefreshToken"), mock.Anything).
			Run(func(args mock.Arguments) {
				capturedRefreshToken = args.Get(1).(*internal_credentials.InternalCredentials)
			}).
			Return(nil)

		var capturedCreds *external_credentials.ExternalCredentials
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.ExternalCredentials"), mock.Anything).
			Run(func(args mock.Arguments) {
				capturedCreds = args.Get(1).(*external_credentials.ExternalCredentials)
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
		assert.Equal(t, command.CredentialsDto.AccessTokenExpiresAtUtc, capturedCreds.AccessTokenExpiresAtUtc())
		assert.Equal(t, command.CredentialsDto.RefreshTokenExpiresAtUtc, capturedCreds.RefreshTokenExpiresAtUtc())
		assert.Equal(t, command.CredentialsDto.ExternalId, capturedCreds.ExternalId())
		assert.Equal(t, command.CredentialsDto.MusicProvider, capturedCreds.MusicProvider())
	})
}
