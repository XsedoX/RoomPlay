package logout_user_command_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_internal_credentials_repository.MockInternalCredentialsRepository,
	*mock_unit_of_work.MockUnitOfWork,
) {
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockRefreshTokenRepository := new(mock_internal_credentials_repository.MockInternalCredentialsRepository)

	defer func() {
		mockUoW.AssertExpectations(t)
		mockRefreshTokenRepository.AssertExpectations(t)
	}()

	return mockRefreshTokenRepository, mockUoW
}

func TestLogoutUserCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenRetireAllTokensByUserIdFails", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		command := &logout_user_command.LogoutUserCommand{
			DeviceId: nil,
			UserId:   user_id.NewUserId(),
		}
		errCode := "LogoutUserCommandHandler.RetireAllTokensByUserId"
		errOfRepository := errors.New("retire all tokens by userId and deviceId failed")
		mockRefreshTokenRepository.On("RetireAllTokensByUserId", mock.Anything, &command.UserId, mock.Anything).Return(errOfRepository)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(err, &parsedErr))
		assert.Equal(t, parsedErr.Code, errCode)
		assert.Equal(t, parsedErr.Err, errOfRepository)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unexpected)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 1)
	})

	t.Run("ShouldReturnErrorWhenRetireTokenByUserIdAndDeviceIdFails", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		deviceId := device_id.NewDeviceId()
		command := &logout_user_command.LogoutUserCommand{
			DeviceId: &deviceId,
			UserId:   user_id.NewUserId(),
		}
		errCode := "LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId"
		errOfRepository := errors.New("retire all tokens by userId and deviceId failed")
		userSession := user_session.NewUserSession(command.UserId, deviceId)
		mockRefreshTokenRepository.On("RetireTokenByUserSession", mock.Anything, *userSession, mock.Anything).Return(errOfRepository)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		var parsedErr *application_error.ApplicationError
		assert.True(t, errors.As(err, &parsedErr))
		assert.Equal(t, parsedErr.Code, errCode)
		assert.Equal(t, parsedErr.Err, errOfRepository)
		assert.Equal(t, parsedErr.ErrorType, application_error_type.Unexpected)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 0)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireTokenByUserSession", 1)
	})
	t.Run("ShouldReturnSuccessWhenDeviceIdNil", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		command := &logout_user_command.LogoutUserCommand{
			DeviceId: nil,
			UserId:   user_id.NewUserId(),
		}
		mockRefreshTokenRepository.On("RetireAllTokensByUserId", mock.Anything, &command.UserId, mock.Anything).Return(nil)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.NoError(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireTokenByUserIdAndDeviceId", 0)
	})
	t.Run("ShouldReturnSuccessWhenDeviceIdNotNil", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		deviceId := device_id.NewDeviceId()
		command := &logout_user_command.LogoutUserCommand{
			DeviceId: &deviceId,
			UserId:   user_id.NewUserId(),
		}
		userSession := user_session.NewUserSession(command.UserId, deviceId)
		mockRefreshTokenRepository.On("RetireTokenByUserSession", mock.Anything, *userSession, mock.Anything).Return(nil)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.NoError(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 0)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireTokenByUserSession", 1)
	})
}
