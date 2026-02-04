package create_room_command_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_encrypter.MockEncrypter,
	*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*user_id.UserId,
	context.Context,
) {
	mockRoomRepo := new(mock_room_repository.MockRoomRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockEncrypter := new(mock_encrypter.MockEncrypter)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
	}()
	return mockEncrypter, mockRoomRepo, mockUoW, &userId, ctx
}

func TestCreateRoomCommandHandler(t *testing.T) {
	t.Run("ShouldCreateRoomSuccessfullyWhenCommandIsValid", func(t *testing.T) {
		// Arrange
		mockEncrypter,
			mockRoomRepo,
			mockUoW,
			_,
			ctx := setupMocks(t)

		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		command := &create_room_command.CreateRoomCommand{
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
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		assert.NoError(t, err)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 1)
	})

	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockEncrypter, mockRoomRepo, mockUoW, _, _ := setupMocks(t)
		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		command := &create_room_command.CreateRoomCommand{
			RoomName:     "Test Room",
			RoomPassword: "password123",
		}

		// Act
		err := handler.Handle(context.Background(), command)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 0)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
	})

	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		// Arrange
		mockEncrypter, mockRoomRepo, mockUoW, _, ctx := setupMocks(t)
		handler := NewCreateRoomCommandHandler(mockRoomRepo, mockUoW, mockEncrypter)
		command := &create_room_command.CreateRoomCommand{
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
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr), "error should be a CustomError")
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "CreateRoom", 1)
	})
}
