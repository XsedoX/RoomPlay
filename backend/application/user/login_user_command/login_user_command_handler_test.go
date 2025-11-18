package login_user_command

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
	"xsedox.com/main/tests"
	"xsedox.com/main/tests/infrustructure/authentication"
	persistance2 "xsedox.com/main/tests/infrustructure/persistance"
)

func TestLoginUserCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenExternalIdMethodFails", func(t *testing.T) {
		// Arrange
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)
		var loginUserCommand = &LoginUserCommand{}
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
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
	})
	t.Run("ShouldReturnErrorWhenUserUpdateFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)
		deviceId := user.DeviceId(uuid.New())

		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		userId := user.Id(uuid.New())
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockedUserFromDb := user.HydrateUser(
			userId,
			loginUserCommand.ExternalId,
			faker.FirstName(),
			faker.LastName(),
			nil,
			nil,
			devices,
			nil,
		)
		mockUserRepository.
			On("GetUserByExternalId", mock.Anything, loginUserCommand.ExternalId, mock.Anything).
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
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 2)
	})
	t.Run("ShouldLoginWithNewDeviceWhenDeviceIdIsNil", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, _ := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		loginUserCommand := &LoginUserCommand{
			Name:       faker.FirstName(),
			DeviceDto:  DeviceDto{DeviceId: nil, DeviceType: user.Mobile},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		// Ensure mockUser is hydrated with at least one device
		devices := []user.Device{
			*user.HydrateDevice(
				user.DeviceId(uuid.New()),
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockUser := user.HydrateUser(
			userId,
			loginUserCommand.ExternalId,
			loginUserCommand.Name,
			loginUserCommand.Surname,
			nil,
			nil,
			devices,
			nil,
		)
		mockUserRepository.
			On("GetUserByExternalId", context.Background(), loginUserCommand.ExternalId, mock.Anything).
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
		assert.NotEqual(t, user.DeviceId{}, resp.DeviceId)
		updateUserCall := mockUserRepository.Calls[1]
		updatedUser := updateUserCall.Arguments.Get(1).(*user.User)
		assert.Len(t, updatedUser.Devices(), 2)
		assert.Equal(t, updatedUser.FullName().Name(), loginUserCommand.Name)
		assert.Equal(t, updatedUser.FullName().Surname(), loginUserCommand.Surname)
		assert.Equal(t, "refresh-token-value", resp.RefreshToken)
		assert.Equal(t, "access-token-value", resp.AccessToken)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 4)
	})
	t.Run("ShouldReloginWithKnownDevice", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		handler := NewLoginUserCommandHandler(mockUoW,
			mockUserRepository,
			mockEncrypter,
			mockJwtProvider,
			mockRefreshTokenRepository,
			mockExternalCredentialsRepository,
		)

		userId := user.Id(uuid.New())
		deviceId := user.DeviceId(uuid.New())
		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              uuid.New().String(),
				RefreshToken:             uuid.New().String(),
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Offline,
				time.Now().Add(-1*time.Hour).UTC(),
			),
		}
		mockUserRepository.
			On("GetUserByExternalId", context.Background(), loginUserCommand.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.ExternalId,
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
		assert.Equal(t, updatedUser.Devices()[0].State(), user.Online)

		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserByExternalId", 1)
		mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 4)
	})
	t.Run("ShouldCallGrantWithCorrectValues", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		userId := user.Id(uuid.New())
		deviceId := user.DeviceId(uuid.New())
		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.
			On("GetUserByExternalId", mock.Anything, loginUserCommand.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.ExternalId,
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

		var capturedCreds *credentials.External
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.External"), mock.Anything).Run(func(args mock.Arguments) {
			capturedCreds = args.Get(1).(*credentials.External)
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
		assert.Equal(t, strings.Split(loginUserCommand.CredentialsDto.Scopes, " "), capturedCreds.Scopes())
		assert.Equal(t, loginUserCommand.CredentialsDto.AccessTokenExpiresAtUtc, capturedCreds.AccessTokenExpiresAtUtc())
		assert.Equal(t, loginUserCommand.CredentialsDto.RefreshTokenExpiresAtUtc, capturedCreds.RefreshTokenExpiresAtUtc())
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("ShouldReturnErrorWhenAssignNewTokenFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		userId := user.Id(uuid.New())
		deviceId := user.DeviceId(uuid.New())
		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.ExternalId,
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
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, "LoginUserCommandHandler.AssignNewToken", customErr.Code)
		assert.ErrorIs(t, customErr.Err, assignErr)
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("ShouldReturnErrorWhenGrantFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		userId := user.Id(uuid.New())
		deviceId := user.DeviceId(uuid.New())
		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.ExternalId,
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
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.External"), mock.Anything).Return(grantErr)

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
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, "LoginUserCommandHandler.Grant", customErr.Code)
		assert.ErrorIs(t, customErr.Err, grantErr)
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockJwtProvider.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("ShouldReturnErrorWhenGenerateTokenFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter := new(persistance2.MockEncrypter)
		mockJwtProvider := new(authentication.MockJwtProvider)
		mockRefreshTokenRepository := new(persistance2.MockRefreshTokenRepository)
		mockExternalCredentialsRepository := new(persistance2.MockExternalCredentialsRepository)
		userId := user.Id(uuid.New())
		deviceId := user.DeviceId(uuid.New())
		loginUserCommand := &LoginUserCommand{
			Name: faker.FirstName(),
			DeviceDto: DeviceDto{
				DeviceId:   &deviceId,
				DeviceType: user.Mobile,
			},
			ExternalId: "SampleExternalId",
			Surname:    faker.LastName(),
			CredentialsDto: CredentialsDto{
				AccessToken:              "access-token-xyz",
				RefreshToken:             "refresh-token-xyz",
				Scopes:                   "openid offline_access",
				AccessTokenExpiresAtUtc:  time.Now().Add(time.Hour).UTC(),
				RefreshTokenExpiresAtUtc: time.Now().Add(time.Hour * 24 * 7).UTC(),
				IssuedAt:                 time.Now().UTC(),
			},
		}
		devices := []user.Device{
			*user.HydrateDevice(
				deviceId,
				"Test Device",
				user.Mobile,
				false,
				user.Online,
				time.Now().UTC(),
			),
		}
		mockUserRepository.On("GetUserByExternalId", mock.Anything, loginUserCommand.ExternalId, mock.Anything).
			Return(user.HydrateUser(
				userId,
				loginUserCommand.ExternalId,
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
		mockExternalCredentialsRepository.On("Grant", mock.Anything, mock.AnythingOfType("*credentials.External"), mock.Anything).Return(nil)

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
