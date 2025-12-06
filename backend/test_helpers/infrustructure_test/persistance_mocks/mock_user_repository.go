package persistance_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Add(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {
	args := m.Called(ctx, user, queryer)
	return args.Error(0)
}

func (m *MockUserRepository) CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer contracts.IQueryer) bool {
	args := m.Called(ctx, externalId, queryer)
	return args.Bool(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {
	args := m.Called(ctx, user, queryer)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByExternalId(ctx context.Context, externalId string, queryer contracts.IQueryer) (*user.User, error) {
	args := m.Called(ctx, externalId, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id user.Id, queryer contracts.IQueryer) (*user.User, error) {
	args := m.Called(ctx, id, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
