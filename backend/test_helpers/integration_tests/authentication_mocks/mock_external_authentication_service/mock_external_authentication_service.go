package mock_external_authentication_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/oidc_authenticate_user_service_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockExternalAuthenticationService struct {
	mock.Mock
}

func (m *MockExternalAuthenticationService) AuthenticateWithExternalProvider(ctx context.Context, code string, deviceId *device_id.DeviceId, deviceType *device_type.DeviceType, provider music_provider.MusicProvider) (*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto, error) {
	args := m.Called(ctx, code, deviceId, deviceType, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto), args.Error(1)
}

func (m *MockExternalAuthenticationService) RefreshAccessTokenWithExternalProvider(ctx context.Context, userId user_id.UserId) (*token.Token, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*token.Token), args.Error(1)
}
