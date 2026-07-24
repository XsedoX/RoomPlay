package mock_cache

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/stretchr/testify/mock"
)

type MockCache[T any] struct {
	mock.Mock
}

func (m *MockCache[T]) Set(key string, value T, ctx context.Context, queryer i_queryer.IQueryer) error {
	args := m.Called(key, value, ctx, queryer)
	return args.Error(0)
}

func (m *MockCache[T]) Remove(key string, ctx context.Context, queryer i_queryer.IQueryer) error {
	args := m.Called(key, ctx, queryer)
	return args.Error(0)
}

func (m *MockCache[T]) Get(key string, ctx context.Context, queryer i_queryer.IQueryer) (T, error) {
	args := m.Called(key, ctx, queryer)
	if args.Get(0) == nil {
		var zero T
		return zero, args.Error(1)
	}
	return args.Get(0).(T), args.Error(1)
}

func (m *MockCache[T]) GetExact(key string, ctx context.Context, queryer i_queryer.IQueryer) (T, error) {
	args := m.Called(key, ctx, queryer)
	if args.Get(0) == nil {
		var zero T
		return zero, args.Error(1)
	}
	return args.Get(0).(T), args.Error(1)
}
