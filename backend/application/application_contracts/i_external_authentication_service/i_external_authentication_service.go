package i_external_authentication_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/oidc_authenticate_user_service_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IExternalAuthenticationService interface {
	AuthenticateWithExternalProvider(ctx context.Context, code string, deviceId *device_id.DeviceId, deviceType *device_type.DeviceType, provider music_provider.MusicProvider) (*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto, error)
	RefreshAccessTokenWithExternalProvider(ctx context.Context, userId user_id.UserId) (*token.Token, error)
}
