package contracts

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/dtos"
	"xsedox.com/main/domain/shared"
)

type IOidcAuthenticationService interface {
	AuthenticateWithGoogle(ctx context.Context, code string, deviceId *shared.DeviceId) (*dtos.OidcAuthenticateUserServiceDto, *applicationErrors.ApplicationError)
}
