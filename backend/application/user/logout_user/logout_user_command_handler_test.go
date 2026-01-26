package logout_user

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*persistance_mocks.MockRefreshTokenRepository,
	*persistance_mocks.MockUnitOfWork,
) {
	mockUoW := new(persistance_mocks.MockUnitOfWork)
	mockRefreshTokenRepository := new(persistance_mocks.MockRefreshTokenRepository)

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
		command := &LogoutUserCommand{
			DeviceId: nil,
			UserId:   user.Id(uuid.New()),
		}
		errCode := "LogoutUserCommandHandler.RetireAllTokensByUserId"
		errOfRepository := errors.New("retire all tokens by userId and deviceId failed")
		mockRefreshTokenRepository.On("RetireAllTokensByUserId", mock.Anything, &command.UserId, mock.Anything).Return(errOfRepository)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		var parsedErr *customerrors.CustomError
		assert.True(t, errors.As(err, &parsedErr))
		assert.Equal(t, parsedErr.Code, errCode)
		assert.Equal(t, parsedErr.Err, errOfRepository)
		assert.Equal(t, parsedErr.ErrorType, customerrors.Unexpected)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 1)
	})

	t.Run("ShouldReturnErrorWhenRetireTokenByUserIdAndDeviceIdFails", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		deviceId := user.DeviceId(uuid.New())
		command := &LogoutUserCommand{
			DeviceId: &deviceId,
			UserId:   user.Id(uuid.New()),
		}
		errCode := "LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId"
		errOfRepository := errors.New("retire all tokens by userId and deviceId failed")
		mockRefreshTokenRepository.On("RetireTokenByUserIdAndDeviceId", mock.Anything, &command.UserId, &deviceId, mock.Anything).Return(errOfRepository)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		var parsedErr *customerrors.CustomError
		assert.True(t, errors.As(err, &parsedErr))
		assert.Equal(t, parsedErr.Code, errCode)
		assert.Equal(t, parsedErr.Err, errOfRepository)
		assert.Equal(t, parsedErr.ErrorType, customerrors.Unexpected)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 0)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireTokenByUserIdAndDeviceId", 1)
	})
	t.Run("ShouldReturnSuccessWhenDeviceIdNil", func(t *testing.T) {
		mockRefreshTokenRepository, mockUnitOfWork := setupMocks(t)

		mockUnitOfWork.On("GetQueryer").Return(nil)
		command := &LogoutUserCommand{
			DeviceId: nil,
			UserId:   user.Id(uuid.New()),
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
		deviceId := user.DeviceId(uuid.New())
		command := &LogoutUserCommand{
			DeviceId: &deviceId,
			UserId:   user.Id(uuid.New()),
		}
		mockRefreshTokenRepository.On("RetireTokenByUserIdAndDeviceId", mock.Anything, &command.UserId, &deviceId, mock.Anything).Return(nil)
		handler := NewLogoutUserCommandHandler(mockRefreshTokenRepository, mockUnitOfWork)

		err := handler.Handle(context.Background(), command)

		assert.NoError(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireAllTokensByUserId", 0)
		mockRefreshTokenRepository.AssertNumberOfCalls(t, "RetireTokenByUserIdAndDeviceId", 1)
	})
}
