package get_user_query_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_user_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_user_repository.MockUserRepository,
	*mock_unit_of_work.MockUnitOfWork,
	user_id.UserId,
	context.Context,
) {
	mockUserRepository := new(mock_user_repository.MockUserRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
	}()

	return mockUserRepository, mockUoW, userId, ctx
}

func TestGetUserQueryHandler(t *testing.T) {
	t.Run("ShouldReturnUserSuccess", func(t *testing.T) {
		mockUserRepository, mockUoW, userId, ctx := setupMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		now := time.Now().UTC().Truncate(time.Second)
		deviceID1 := device_id.NewDeviceId()
		deviceID2 := device_id.NewDeviceId()
		userRole := user_role.Member
		devices := []device.Device{
			*device.HydrateDevice(deviceID1,
				faker.Name(),
				device_type.Desktop,
				true,
				device_state.Online,
				now.Add(-5*time.Minute)),
			*device.HydrateDevice(deviceID2,
				faker.Name(),
				device_type.Mobile,
				false,
				device_state.Offline,
				now.Add(-10*time.Minute)),
		}
		userToBeReturned := user.HydrateUser(userId,
			faker.FirstName(),
			faker.LastName(),
			&userRole,
			nil,
			devices,
			nil)
		mockUserRepository.
			On("GetUserById", ctx, userId, mock.Anything).
			Return(userToBeReturned, nil)
		handler := NewGetUserQueryHandler(mockUoW, mockUserRepository)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		assert.Equal(t, resp.Name, userToBeReturned.FullName().Name())
		assert.Equal(t, resp.Surname, userToBeReturned.FullName().Surname())
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockUserRepository, mockUoW, _, _ := setupMocks(t)
		handler := NewGetUserQueryHandler(mockUoW, mockUserRepository)
		// Act
		userObj, err := handler.Handle(context.Background())
		// Assert
		assert.Error(t, err)
		assert.Nil(t, userObj)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 0)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
	})
	t.Run("ShouldReturnErrorWhenUserRepositoryFails", func(t *testing.T) {
		// Arrange
		mockUserRepository, mockUoW, userId, cont := setupMocks(t)
		repoErr := errors.New("database error")
		errorCode := "NewGetUserQueryHandler.GetUserById"
		mockUoW.On("GetQueryer").Return(nil)
		mockUserRepository.On("GetUserById", cont, userId, mock.Anything).Return(nil, repoErr)
		handler := NewGetUserQueryHandler(mockUoW, mockUserRepository)
		// Act
		userObj, err := handler.Handle(cont)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, userObj)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr), "error should be a CustomError")
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
	})
}
