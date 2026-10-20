package mock_short_lived_cache

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockShortLivedCache[T any] struct {
	mock.Mock
}

func (m *MockShortLivedCache[T]) Set(key string, value T, ctx context.Context) error {
	args := m.Called(key, value, ctx)
	return args.Error(0)
}

func (m *MockShortLivedCache[T]) Get(key string, ctx context.Context) (T, error) {
	args := m.Called(key, ctx)
	if args.Get(0) == nil {
		var zero T
		return zero, args.Error(1)
	}
	return args.Get(0).(T), args.Error(1)
}
