package get_user_room_membership_query

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

func TestGetUserRoomMembershipQueryHandler(t *testing.T) {
	t.Run("ShouldReturnBoolSuccess", func(t *testing.T) {
		mockRoomRepo := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())
		mockRoomRepo.
			On("GetRoomByUserId", ctx, userId, mock.Anything).
			Return(nil, nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepo.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		assert.Equal(t, true, *resp)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)

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
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository := new(persistance_mocks.MockRoomRepository)
		mockUoW := new(persistance_mocks.MockUnitOfWork)
		userId, ctx := test_helpers.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "GetUserRoomMembershipQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, userId, mock.Anything).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
