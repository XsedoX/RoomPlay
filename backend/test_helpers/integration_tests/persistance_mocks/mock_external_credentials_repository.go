package persistance_mocks

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockExternalCredentialsRepository struct {
	mock.Mock
}

func (m *MockExternalCredentialsRepository) Grant(ctx context.Context, credentials *credentials.ExternalCredentials, queryer application_contracts.IQueryer) error {
	args := m.Called(ctx, credentials, queryer)
	return args.Error(0)
}

func (m *MockExternalCredentialsRepository) GetAccessTokenByUserId(ctx context.Context, userId user.Id, queryer application_contracts.IQueryer) (string, error) {
	args := m.Called(ctx, userId, queryer)
	return args.String(0), args.Error(1)
}
