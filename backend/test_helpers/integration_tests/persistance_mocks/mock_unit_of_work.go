package persistance_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/XsedoX/RoomPlay/application/contracts"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) ExecuteTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (m *MockUnitOfWork) ExecuteRead(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (m *MockUnitOfWork) GetQueryer() contracts.IQueryer {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(contracts.IQueryer)
}
