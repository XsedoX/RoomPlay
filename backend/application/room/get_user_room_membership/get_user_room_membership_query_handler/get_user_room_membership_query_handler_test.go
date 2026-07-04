package get_user_room_membership_query_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks(t *testing.T) (*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
	user_id.UserId,
	context.Context,
) {
	mockRoomRepo := new(mock_room_repository.MockRoomRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertExpectations(t)
	}()

	return mockRoomRepo, mockUoW, userId, ctx
}

func TestGetUserRoomMembershipQueryHandler(t *testing.T) {
	t.Run("ShouldReturnBoolSuccess", func(t *testing.T) {
		mockRoomRepo, mockUoW, userId, ctx := setupMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		mockRoomRepo.
			On("GetRoomByUserId", ctx, userId, mock.Anything).
			Return(nil, nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		assert.Equal(t, true, *resp)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(mock_room_repository.MockRoomRepository)
		mockUoW := new(mock_unit_of_work.MockUnitOfWork)

		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		// Act
		resp, err := handler.Handle(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepo.AssertNumberOfCalls(t, "GetRoomByUserId", 0)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setupMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "GetUserRoomMembershipQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, userId, mock.Anything).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
