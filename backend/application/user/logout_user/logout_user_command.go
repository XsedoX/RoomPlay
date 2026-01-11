package logout_user

import (
	"github.com/XsedoX/RoomPlay/domain/user"
)

type LogoutUserCommand struct {
	DeviceId *user.DeviceId
	UserId   user.Id
}
