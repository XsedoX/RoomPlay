package persistance

import (
	"context"

	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
)

type MockExternalCredentialsRepository struct {
	mock.Mock
}

func (m *MockExternalCredentialsRepository) Grant(ctx context.Context, credentials *credentials.External, queryer contracts.IQueryer) error {
	args := m.Called(ctx, credentials, queryer)
	return args.Error(0)
}
