package get_user_room_membership_query_handler

import (
	"context"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
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
			On("GetUserMembership", ctx, userId, mock.Anything).
			Return(new(room_id.NewRoomId()), nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "GetUserMembership", 1)
		assert.NotNil(t, resp)
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
		mockRoomRepo.AssertNumberOfCalls(t, "GetUserMembership", 0)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError("GetUserRoomMembershipQueryHandler.Handle"), err)
	})
}
