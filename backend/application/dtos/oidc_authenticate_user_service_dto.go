package dtos

import (
	"xsedox.com/main/domain/user"
)

type OidcAuthenticateUserServiceDto struct {
	RefreshToken string        `json:"refresh_token"`
	AccessToken  string        `json:"access_token"`
	DeviceId     user.DeviceId `json:"device_id"`
}
