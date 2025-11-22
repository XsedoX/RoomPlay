package persistance_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer contracts.IQueryer) error {
	args := m.Called(ctx, refreshToken, queryer)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetTokenByValue(ctx context.Context, value string, queryer contracts.IQueryer) (*credentials.RefreshToken, error) {
	args := m.Called(ctx, value, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*credentials.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user.Id, deviceId *user.DeviceId, queryer contracts.IQueryer) error {
	args := m.Called(ctx, userId, deviceId, queryer)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RetireAllTokensByUserId(ctx context.Context, userId *user.Id, queryer contracts.IQueryer) error {
	args := m.Called(ctx, userId, queryer)
	return args.Error(0)
}
