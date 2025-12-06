package logout_user

import (
	"xsedox.com/main/domain/user"
)

type LogoutUserCommand struct {
	DeviceId *user.DeviceId
	UserId   user.Id
}
