package logout_user_command

import (
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type LogoutUserCommand struct {
	DeviceId *device_id.DeviceId
	UserId   user_id.UserId
}
