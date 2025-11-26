package create_room_command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/test_helpers"
	"xsedox.com/main/test_helpers/infrustructure_test/persistance_mocks"
)

func TestCreateRoomCommandHandler(t *testing.T) {
	t.Run("ShouldCreateRoomSuccessfullyWhenCommandIsValid", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		_, ctx := test_helpers.AddUserIdToContext(context.Background())
		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		command := &CreateRoomCommand{
			RoomName:     "Test Room",
			RoomPassword: "password123",
		}
		encryptionKey := "some-random-key"
		mockUoW.On("GetQueryer").Return(nil)
		mockEncrypter.On("NewEncryptionKey").Return([]byte(encryptionKey))
		mockRoomRepo.On("CreateRoom", ctx, mock.AnythingOfType("*room.Room"), mock.Anything).Return(nil)

		// Act
		err := handler.Handle(ctx, command)

		// Assert
		assert.NoError(t, err)
		mockRoomRepo.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 1)
	})

	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		command := &CreateRoomCommand{
			RoomName:     "Test Room",
			RoomPassword: "password123",
		}

		// Act
		err := handler.Handle(context.Background(), command)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 0)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
	})

	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockEncrypter := new(persistance_mocks.MockEncrypter)
		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		_, ctx := test_helpers.AddUserIdToContext(context.Background())
		command := &CreateRoomCommand{
			RoomName:     "Test Room",
			RoomPassword: "password123",
		}
		repoErr := errors.New("database error")
		errorCode := "CreateRoomCommandHandler.CreateRoom"
		encryptionKey := "some-random-key"
		mockEncrypter.On("NewEncryptionKey").Return([]byte(encryptionKey))
		mockUoW.On("GetQueryer").Return(nil)
		mockRoomRepo.On("CreateRoom", ctx, mock.AnythingOfType("*room.Room"), mock.Anything).Return(repoErr)

		// Act
		err := handler.Handle(ctx, command)

		// Assert
		assert.Error(t, err)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr), "error should be a CustomError")
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 1)
	})
}
