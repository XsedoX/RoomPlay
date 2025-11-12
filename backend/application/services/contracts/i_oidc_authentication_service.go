package contracts

import (
	"context"

	"xsedox.com/main/application/dtos"
	"xsedox.com/main/domain/user"
)

type IOidcAuthenticationService interface {
	AuthenticateWithGoogle(ctx context.Context, code string, deviceId *user.DeviceId, deviceType *user.DeviceType) (*dtos.OidcAuthenticateUserServiceDto, error)
}
