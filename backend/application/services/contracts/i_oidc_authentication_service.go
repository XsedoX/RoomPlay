package contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IOidcAuthenticationService interface {
	AuthenticateWithGoogle(ctx context.Context, code string, deviceId *user.DeviceId, deviceType *user.DeviceType) (*dtos.OidcAuthenticateUserServiceDto, error)
}
