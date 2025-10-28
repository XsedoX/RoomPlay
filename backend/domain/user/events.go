package user

import (
	"xsedox.com/main/domain/shared"
)

type LoggedInEvent struct {
	UserId   shared.UserId
	UserRole Role
}
