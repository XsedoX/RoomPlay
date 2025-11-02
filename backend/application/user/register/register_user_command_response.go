package register

import "xsedox.com/main/domain/shared"

type UserCommandResponse struct {
	RefreshToken string          `json:"refreshToken"`
	AccessToken  string          `json:"accessToken"`
	DeviceId     shared.DeviceId `json:"deviceId"`
}
