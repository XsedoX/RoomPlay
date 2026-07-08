package leave_room_command_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*user_id.UserId,
	context.Context,
) {
	mockRoomRepository := new(mock_room_repository.MockRoomRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
	}()

	return mockRoomRepository, mockUoW, &userId, ctx
}

func TestLeaveRoomCommandHandler(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setupMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		command := &leave_room_command.LeaveRoomCommand{}
		mockRoomRepository.On("LeaveRoom", ctx, *userId, mock.Anything).Return(nil)

		err := handler.Handle(ctx, command)

		assert.NoError(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 1)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepository, mockUoW, _, _ := setupMocks(t)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		command := &leave_room_command.LeaveRoomCommand{}

		// Act
		err := handler.Handle(context.Background(), command)

		// Assert
		assert.Error(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 0)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenUserRepositoryFails", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setupMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockRoomRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "LeaveRoomCommandHandler.LeaveRoom"
		mockRoomRepository.On("LeaveRoom", ctx, *userId, mock.Anything).Return(repoErr)
		command := &leave_room_command.LeaveRoomCommand{}

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		var customErr *application_error.ApplicationError
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "LeaveRoom", 1)
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
