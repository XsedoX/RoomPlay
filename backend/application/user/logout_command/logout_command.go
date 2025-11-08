package logout_command

import (
	"xsedox.com/main/domain/user"
)

type Command struct {
	DeviceId *user.DeviceId
	UserId   user.Id
}
