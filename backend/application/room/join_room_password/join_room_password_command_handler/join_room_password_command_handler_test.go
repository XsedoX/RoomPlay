package join_room_password_command_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
) {
	mockRoomRepository := new(mock_room_repository.MockRoomRepository)
	mockUnitOfWork := new(mock_unit_of_work.MockUnitOfWork)
	defer func() {
		mockRoomRepository.AssertExpectations(t)
		mockUnitOfWork.AssertExpectations(t)
	}()

	return mockRoomRepository, mockUnitOfWork
}

func TestJoinRoomPasswordCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenUserIdMissingFromContext", func(t *testing.T) {
		mockRoomRepository, mockUnitOfWork := setupMocks(t)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &join_room_password_command.JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 0)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenGetRoomIdByNameAndPasswordFails", func(t *testing.T) {
		mockRoomRepository, mockUnitOfWork := setupMocks(t)
		mockUnitOfWork.On("GetQueryer").Return(nil)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &join_room_password_command.JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}
		_, ctx := test_helpers.AddUserIdToContext(context.Background())

		repoErr := errors.New("database error")
		errorCode := "JoinRoomPasswordCommandHandler.GetRoomIdByNameAndPassword"
		mockRoomRepository.On("GetRoomIdByNameAndPassword", ctx, command.RoomName, command.RoomPassword, mock.Anything).
			Return(nil, repoErr)

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 0)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
	t.Run("ShouldReturnErrorWhenJoinRoomByIdFails", func(t *testing.T) {
		mockRoomRepository, mockUnitOfWork := setupMocks(t)
		mockUnitOfWork.On("GetQueryer").Return(nil)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &join_room_password_command.JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())

		repoErr := errors.New("database error")
		errorCode := "JoinRoomPasswordCommandHandler.JoinRoomById"
		roomId := room_id.NewRoomId()
		mockRoomRepository.On("GetRoomIdByNameAndPassword", ctx, command.RoomName, command.RoomPassword, mock.Anything).
			Return(&roomId, nil).
			On("JoinRoomById", ctx, userId, roomId, mock.Anything).
			Return(repoErr).
			Once()

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 2)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 1)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockRoomRepository, mockUnitOfWork := setupMocks(t)
		mockUnitOfWork.On("GetQueryer").Return(nil).Twice()
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &join_room_password_command.JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())

		roomId := room_id.NewRoomId()
		mockRoomRepository.On("GetRoomIdByNameAndPassword", ctx, command.RoomName, command.RoomPassword, mock.Anything).
			Return(&roomId, nil).
			Once().
			On("JoinRoomById", ctx, userId, roomId, mock.Anything).
			Return(nil).
			Once()

		err := handler.Handle(ctx, command)

		assert.NoError(t, err)
	})
}
