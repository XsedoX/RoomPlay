package persistance_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
)

type MockExternalCredentialsRepository struct {
	mock.Mock
}

func (m *MockExternalCredentialsRepository) Grant(ctx context.Context, credentials *credentials.External, queryer contracts.IQueryer) error {
	args := m.Called(ctx, credentials, queryer)
	return args.Error(0)
}
