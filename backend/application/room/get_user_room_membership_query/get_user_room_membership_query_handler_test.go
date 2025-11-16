package get_user_room_membership_query

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

func TestGetUserRoomMembershipQueryHandler(t *testing.T) {
	t.Run("ShouldReturnBoolSuccess", func(t *testing.T) {
		mockRoomRepo := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		mockUoW.On("GetQueryer").Return(nil)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockRoomRepo.
			On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).
			Return(nil, nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertExpectations(t)
		assert.Equal(t, true, *resp)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)

		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepo, mockUoW)

		// Act
		resp, err := handler.Handle(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockUoW.AssertExpectations(t)
		mockRoomRepo.AssertExpectations(t)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetUserRoomMembershipQueryHandler(mockRoomRepository, mockUoW)
		repoErr := errors.New("database error")
		errorCode := "GetUserRoomMembershipQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		mockUoW.AssertExpectations(t)
		mockRoomRepository.AssertExpectations(t)
		assert.ErrorIs(t, customErr.Err, repoErr)
	})
}
