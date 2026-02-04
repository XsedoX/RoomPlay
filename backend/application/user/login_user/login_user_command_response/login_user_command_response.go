package login_user_command_response

import (
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
)

type LoginUserCommandResponse struct {
	RefreshToken string             `json:"refreshToken"`
	AccessToken  string             `json:"accessToken"`
	DeviceId     device_id.DeviceId `json:"deviceId"`
}
