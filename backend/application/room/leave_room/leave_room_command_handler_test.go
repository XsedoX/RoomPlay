package leave_room

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/test_helpers"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaveRoomCommandHandler(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		command := &LeaveRoomCommand{}
		mockRoomRepository.On("LeaveRoom", ctx, userId, mock.Anything).Return(nil)

		err := handler.Handle(ctx, command)

		assert.NoError(t, err)
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 1)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		command := &LeaveRoomCommand{}

		// Act
		err := handler.Handle(context.Background(), command)

		// Assert
		assert.Error(t, err)
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 0)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenUserRepositoryFails", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "LeaveRoomCommandHandler.LeaveRoom"
		mockRoomRepository.On("LeaveRoom", ctx, userId, mock.Anything).Return(repoErr)
		command := &LeaveRoomCommand{}

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		var customErr *customerrors.CustomError
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 1)
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
