package mock_internal_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockInternalCredentialsRepository struct {
	mock.Mock
}

func (m *MockInternalCredentialsRepository) AssignNewToken(ctx context.Context, refreshToken *internal_credentials.InternalCredentials, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, refreshToken, queryer)
	return args.Error(0)
}

func (m *MockInternalCredentialsRepository) GetTokenByValue(ctx context.Context, value string, queryer i_queryer.IQueryer) (*internal_credentials.InternalCredentials, error) {
	args := m.Called(ctx, value, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*internal_credentials.InternalCredentials), args.Error(1)
}

func (m *MockInternalCredentialsRepository) RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user_id.UserId, deviceId *device_id.DeviceId, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, userId, deviceId, queryer)
	return args.Error(0)
}

func (m *MockInternalCredentialsRepository) RetireAllTokensByUserId(ctx context.Context, userId *user_id.UserId, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, userId, queryer)
	return args.Error(0)
}
