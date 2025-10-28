package dtos

import "xsedox.com/main/domain/shared"

type OidcAuthenticateUserServiceDto struct {
	RefreshToken string          `json:"refresh_token"`
	AccessToken  string          `json:"access_token"`
	DeviceId     shared.DeviceId `json:"device_id"`
}
