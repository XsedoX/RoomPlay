package oidc_authenticate_user_service_dto

import (
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
)

type OidcAuthenticateUserServiceDto struct {
	RefreshToken string             `json:"refresh_token"`
	AccessToken  string             `json:"access_token"`
	DeviceId     device_id.DeviceId `json:"device_id"`
}
