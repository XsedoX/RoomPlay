package mock_user_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Add(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, user, queryer)
	return args.Error(0)
}

func (m *MockUserRepository) CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) bool {
	args := m.Called(ctx, externalId, queryer)
	return args.Bool(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, user, queryer)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) (*user.User, error) {
	args := m.Called(ctx, externalId, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) (*user.User, error) {
	args := m.Called(ctx, id, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
