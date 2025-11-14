package logout_user_command

import (
	"xsedox.com/main/domain/user"
)

type LogoutUserCommand struct {
	DeviceId *user.DeviceId
	UserId   user.Id
}
