package join_room_password

import (
	"context"
	"errors"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/test_helpers"
	"xsedox.com/main/test_helpers/infrustructure_test/persistance_mocks"
)

func TestJoinRoomPasswordCommandHandler(t *testing.T) {
	t.Run("ShouldReturnErrorWhenUserIdMissingFromContext", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUnitOfWork := new(persistance_mocks.MockUnitOfWork)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}

		err := handler.Handle(context.Background(), command)

		assert.Error(t, err)
		mockUnitOfWork.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 0)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenGetRoomIdByNameAndPasswordFails", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUnitOfWork := new(persistance_mocks.MockUnitOfWork)
		mockUnitOfWork.On("GetQueryer").Return(nil)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &JoinRoomPasswordCommand{
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
		mockUnitOfWork.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 0)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
	t.Run("ShouldReturnErrorWhenJoinRoomByIdFails", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUnitOfWork := new(persistance_mocks.MockUnitOfWork)
		mockUnitOfWork.On("GetQueryer").Return(nil)
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())

		repoErr := errors.New("database error")
		errorCode := "JoinRoomPasswordCommandHandler.JoinRoomById"
		roomId := shared.RoomId(uuid.New())
		mockRoomRepository.On("GetRoomIdByNameAndPassword", ctx, command.RoomName, command.RoomPassword, mock.Anything).
			Return(&roomId, nil).
			On("JoinRoomById", ctx, userId, roomId, mock.Anything).
			Return(repoErr).
			Once()

		err := handler.Handle(ctx, command)

		assert.Error(t, err)
		mockUnitOfWork.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUnitOfWork.AssertNumberOfCalls(t, "GetQueryer", 2)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomIdByNameAndPassword", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "JoinRoomById", 1)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUnitOfWork := new(persistance_mocks.MockUnitOfWork)
		mockUnitOfWork.On("GetQueryer").Return(nil).Twice()
		handler := NewJoinRoomPasswordCommandHandler(mockRoomRepository, mockUnitOfWork)
		command := &JoinRoomPasswordCommand{
			RoomName:     faker.Word(),
			RoomPassword: faker.Password(),
		}
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())

		roomId := shared.RoomId(uuid.New())
		mockRoomRepository.On("GetRoomIdByNameAndPassword", ctx, command.RoomName, command.RoomPassword, mock.Anything).
			Return(&roomId, nil).
			Once().
			On("JoinRoomById", ctx, userId, roomId, mock.Anything).
			Return(nil).
			Once()

		err := handler.Handle(ctx, command)

		assert.NoError(t, err)
		mockUnitOfWork.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
	})
}
