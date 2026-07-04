package i_oidc_authentication_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/oidc_authenticate_user_service_dto"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
)

type IOidcAuthenticationService interface {
	AuthenticateWithGoogle(ctx context.Context, code string, deviceId *device_id.DeviceId, deviceType *device_type.DeviceType) (*oidc_authenticate_user_service_dto.OidcAuthenticateUserServiceDto, error)
}
