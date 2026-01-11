package login_user

import (
	"github.com/XsedoX/RoomPlay/domain/user"
)

type LoginUserCommandResponse struct {
	RefreshToken string        `json:"refreshToken"`
	AccessToken  string        `json:"accessToken"`
	DeviceId     user.DeviceId `json:"deviceId"`
}
