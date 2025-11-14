package login_user_command

import (
	"xsedox.com/main/domain/user"
)

type LoginUserCommandResponse struct {
	RefreshToken string        `json:"refreshToken"`
	AccessToken  string        `json:"accessToken"`
	DeviceId     user.DeviceId `json:"deviceId"`
}
