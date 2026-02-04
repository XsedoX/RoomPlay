package login_user_command_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_jwt_provider"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_user_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_user_repository.MockUserRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*mock_encrypter.MockEncrypter,
	*mock_jwt_provider.MockJwtProvider,
	*mock_internal_credentials_repository.MockInternalCredentialsRepository,
	*mock_external_credentials_repository.MockExternalCredentialsRepository,
	user_id.UserId,
	context.Context,
) {
	mockUserRepository := new(mock_user_repository.MockUserRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockEncrypter := new(mock_encrypter.MockEncrypter)
	mockJwtProvider := new(mock_jwt_provider.MockJwtProvider)
	mockRefreshTokenRepository := new(mock_internal_credentials_repository.MockInternalCredentialsRepository)
	mockExternalCredentialsRepository := new(mock_external_credentials_repository.MockExternalCredentialsRepository)

	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	}()

	return mockUserRepository, mockUoW, mockEncrypter, mockJwtProvider, mockRefreshTokenRepository, mockExternalCredentialsRepository, userId, ctx
}

func TestLoginUserCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenExternalIdMethodFails", func(t *testing.T) {
		// Arrange
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)
		loginUserCommand := &login_user_command.LoginUserCommand{}
		repoErr := errors.New("database error")
		errorCode := "LoginUserCommandHandler.GetUserByExternalId"
		mockUserRepository.
			On("GetUserByExternalId", context.Background(), mock.Anything, mock.Anything).
			Return(nil, repoErr)
		// Act
		resp, err := handler.Handle(context.Background(), loginUserCommand)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
	})
	t.Run("ShouldReturnErrorWhenUserUpdateFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)
		deviceId := device_id.NewDeviceId()

		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               "SampleExternalId",
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		userId := user_id.NewUserId()
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockedUserFromDb := user.HydrateUser(
			userId,
			faker.FirstName(),
			faker.LastName(),
			nil,
			nil,
			devices,
			nil,
		)
		mockUserRepository.
			On("GetUserByExternalId", mock.Anything, loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(mockedUserFromDb, nil)
		repoErr := errors.New("update error")
		errorCode := "LoginUserCommandHandler.Update"
		mockUserRepository.
			On("Update", context.Background(), mockedUserFromDb, mock.Anything).
			Return(repoErr)
		// Act
		resp, err := handler.Handle(context.Background(), loginUserCommand)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 2)
	})
	t.Run("ShouldLoginWithNewDeviceWhenDeviceIdIsNil", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			userId,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		loginUserCommand := &login_user_command.LoginUserCommand{
			Name:      faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{DeviceId: nil, DeviceType: device_type.Mobile},
			Surname:   faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				AccessToken:              uuid.New().String(),
				ExternalId:               "SampleExternalId",
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		// Ensure mockUser is hydrated with at least one device
		devices := []device.Device{
			*device.HydrateDevice(
				device_id.NewDeviceId(),
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockUser := user.HydrateUser(
			userId,
			loginUserCommand.Name,
			loginUserCommand.Surname,
			nil,
			nil,
			devices,
			nil,
		)
		mockUserRepository.
			On("GetUserByExternalId", context.Background(), loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(mockUser, nil)
		mockUserRepository.
			On("Update", context.Background(), mock.Anything, mock.Anything).
			Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		mockRefreshTokenRepository.On("AssignNewToken", context.Background(), mock.Anything, mock.Anything).Return(nil)
		mockJwtProvider.On("GenerateToken", mockUser.Id()).Return("access-token-value", nil)
		mockExternalCredentialsRepository.On("Grant", context.Background(), mock.Anything, mock.Anything).Return(nil)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEqual(t, device_id.DeviceId{}, resp.DeviceId)
		updateUserCall := mockUserRepository.Calls[1]
		updatedUser := updateUserCall.Arguments.Get(1).(*user.User)
		assert.Len(t, updatedUser.Devices(), 2)
		assert.Equal(t, updatedUser.FullName().Name(), loginUserCommand.Name)
		assert.Equal(t, updatedUser.FullName().Surname(), loginUserCommand.Surname)
		assert.Equal(t, "refresh-token-value", resp.RefreshToken)
		assert.Equal(t, "access-token-value", resp.AccessToken)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 4)
	})
	t.Run("ShouldReloginWithKnownDevice", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		userId := user_id.NewUserId()
		deviceId := device_id.NewDeviceId()
		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               "SampleExternalId",
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Offline,
				time.Now().Add(-1*time.Hour).UTC(),
			),
		}
		mockUserRepository.
			On("GetUserByExternalId", context.Background(), loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				faker.FirstName(),
				faker.LastName(),
				nil,
				nil,
				devices,
				nil,
			), nil)
		mockUserRepository.
			On("Update", context.Background(), mock.Anything, mock.Anything).
			Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		mockRefreshTokenRepository.On("AssignNewToken", context.Background(), mock.Anything, mock.Anything).Return(nil)
		// Use hydrated user for GenerateToken expectation
		mockJwtProvider.On("GenerateToken", userId).Return("access-token-value", nil)
		mockExternalCredentialsRepository.On("Grant", context.Background(), mock.Anything, mock.Anything).Return(nil)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "refresh-token-value", resp.RefreshToken)
		assert.Equal(t, "access-token-value", resp.AccessToken)
		assert.Equal(t, resp.DeviceId, deviceId)
		updateUserCall := mockUserRepository.Calls[1]
		updatedUser := updateUserCall.Arguments.Get(1).(*user.User)
		assert.Less(t, time.Now().Add(-1*time.Minute), updatedUser.Devices()[0].LastLoggedInUtc())
		assert.Equal(t, updatedUser.Devices()[0].Id(), deviceId)
		assert.Equal(t, updatedUser.Devices()[0].State(), device_state.Online)

		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 4)
	})
	t.Run("ShouldCallGrantWithCorrectValues", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		userId := user_id.NewUserId()
		deviceId := device_id.NewDeviceId()
		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               "SampleExternalId",
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.
			On("GetUserByExternalId", mock.Anything, loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.Name,
				loginUserCommand.Surname,
				nil,
				nil,
				devices,
				nil,
			), nil)
		mockUserRepository.
			On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockJwtProvider.On("GenerateToken", userId).Return("access-token-value", nil)

		var capturedCreds *external_credentials.ExternalCredentials
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.ExternalCredentials"), mock.Anything).Run(func(args mock.Arguments) {
			capturedCreds = args.Get(1).(*external_credentials.ExternalCredentials)
		}).Return(nil)

		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, capturedCreds)
		assert.Equal(t, userId, capturedCreds.Id())
		assert.Equal(t, loginUserCommand.CredentialsDto.AccessToken, capturedCreds.AccessToken())
		assert.Equal(t, loginUserCommand.CredentialsDto.RefreshToken, capturedCreds.RefreshToken())
		assert.Equal(t, loginUserCommand.CredentialsDto.MusicProvider, capturedCreds.MusicProvider())
		assert.Equal(t, loginUserCommand.CredentialsDto.AccessTokenExpiresAtUtc, capturedCreds.AccessTokenExpiresAtUtc())
		assert.Equal(t, loginUserCommand.CredentialsDto.RefreshTokenExpiresAtUtc, capturedCreds.RefreshTokenExpiresAtUtc())
	})
	t.Run("ShouldReturnErrorWhenAssignNewTokenFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		userId := user_id.NewUserId()
		deviceId := device_id.NewDeviceId()
		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               "SampleExternalId",
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.Name,
				loginUserCommand.Surname,
				nil,
				nil,
				devices,
				nil,
			), nil)
		mockUserRepository.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		assignErr := errors.New("assign token error")
		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.Anything, mock.Anything).Return(assignErr)

		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, "LoginUserCommandHandler.AssignNewToken", customErr.Code)
		assert.ErrorIs(t, customErr.Err, assignErr)
	})
	t.Run("ShouldReturnErrorWhenGrantFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		userId := user_id.NewUserId()
		deviceId := device_id.NewDeviceId()
		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				AccessToken:              "access-token-xyz",
				ExternalId:               "SampleExternalId",
				RefreshToken:             "refresh-token-xyz",
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.Name,
				loginUserCommand.Surname,
				nil,
				nil,
				devices,
				nil,
			), nil)
		mockUserRepository.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockJwtProvider.On("GenerateToken", userId).Return("access-token-value", nil)
		grantErr := errors.New("grant error")
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.ExternalCredentials"), mock.Anything).Return(grantErr)

		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, "LoginUserCommandHandler.Grant", customErr.Code)
		assert.ErrorIs(t, customErr.Err, grantErr)
	})
	t.Run("ShouldReturnErrorWhenGenerateTokenFails", func(t *testing.T) {
		mockUserRepository,
			mockUoW,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
			_,
			_ := setupMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		userId := user_id.NewUserId()
		deviceId := device_id.NewDeviceId()
		loginUserCommand := &login_user_command.LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: login_user_command.DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: device_type.Mobile,
			},
			Surname: faker.LastName(),
			CredentialsDto: login_user_command.CredentialsDto{
				ExternalId:               "SampleExternalId",
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				MusicProvider:            music_provider.YouTube,
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []device.Device{
			*device.HydrateDevice(
				deviceId,
				"Test Device",
				device_type.Mobile,
				false,
				device_state.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.CredentialsDto.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.Name,
				loginUserCommand.Surname,
				nil,
				nil,
				devices,
				nil,
			), nil)
		mockUserRepository.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte("refresh-token-value"))
		mockRefreshTokenRepository.On("AssignNewToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		tokenErr := errors.New("token generation error")
		mockJwtProvider.On("GenerateToken", userId).Return("", tokenErr)
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.ExternalCredentials"), mock.Anything).Return(nil)

		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		resp, err := handler.Handle(context.Background(), loginUserCommand)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, tokenErr)
	})
}
