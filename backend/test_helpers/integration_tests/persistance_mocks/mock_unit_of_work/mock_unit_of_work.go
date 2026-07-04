package mock_unit_of_work

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/stretchr/testify/mock"
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

func (m *MockUnitOfWork) GetQueryer() i_queryer.IQueryer {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(i_queryer.IQueryer)
}
