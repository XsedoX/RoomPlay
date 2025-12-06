package get_user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/test_helpers"
	"xsedox.com/main/test_helpers/infrustructure_test/persistance_mocks"
)

func TestGetUserQueryHandler(t *testing.T) {
	t.Run("ShouldReturnUserSuccess", func(t *testing.T) {
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())
		now := time.Now().UTC().Truncate(time.Second)
		deviceID1 := user.DeviceId(uuid.New())
		deviceID2 := user.DeviceId(uuid.New())
		userRole := user.Member
		devices := []user.Device{
			*user.HydrateDevice(deviceID1,
				faker.Name(),
				user.Desktop,
				true,
				user.Online,
				now.Add(-5*time.Minute)),
			*user.HydrateDevice(deviceID2,
				faker.Name(),
				user.Mobile,
				false,
				user.Offline,
				now.Add(-10*time.Minute)),
		}
		userToBeReturned := user.HydrateUser(userId,
			uuid.New().String(),
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
		mockUserRepository.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockUoW.AssertExpectations(t)
		assert.Equal(t, resp.Name, userToBeReturned.FullName().Name())
		assert.Equal(t, resp.Surname, userToBeReturned.FullName().Surname())
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		handler := NewGetUserQueryHandler(mockUoW, mockUserRepository)
		// Act
		userObj, err := handler.Handle(context.Background())
		// Assert
		assert.Error(t, err)
		assert.Nil(t, userObj)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 0)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
	})
	t.Run("ShouldReturnErrorWhenUserRepositoryFails", func(t *testing.T) {
		// Arrange
		mockUserRepository := new(persistance_mocks.MockUserRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		userId, cont := test_helpers.AddUserIdToContext(context.Background())
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
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr), "error should be a CustomError")
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUserRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertNumberOfCalls(t, "GetUserById", 1)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
	})
}
