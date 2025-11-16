package leave_room_command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"xsedox.com/main/application"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/tests"
	persistance2 "xsedox.com/main/tests/infrustructure/persistance"
)

func TestLeaveRoomCommandHandler(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockUserRepository, mockUoW)
		command := &LeaveRoomCommand{}
		mockUserRepository.On("LeaveRoom", ctx, userId, mockUoW.GetQueryer()).Return(nil)

		err := handler.Handle(ctx, command)

		assert.NoError(t, err)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		handler := NewLeaveRoomCommandHandler(mockUserRepository, mockUoW)
		command := &LeaveRoomCommand{}

		// Act
		err := handler.Handle(context.Background(), command)

		// Assert
		assert.Error(t, err)
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockUserRepository := new(persistance2.MockUserRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewLeaveRoomCommandHandler(mockUserRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "LeaveRoomCommandHandler.LeaveRoom"
		mockUserRepository.On("LeaveRoom", ctx, userId, mockUoW.GetQueryer()).Return(repoErr)
		command := &LeaveRoomCommand{}

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		var customErr *custom_errors.CustomError
		mockUoW.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
