package register

import (
	"xsedox.com/main/domain/user"
)

type UserCommandResponse struct {
	RefreshToken string        `json:"refreshToken"`
	AccessToken  string        `json:"accessToken"`
	DeviceId     user.DeviceId `json:"deviceId"`
}
