package login_user_command

import (
	"time"

	"xsedox.com/main/domain/user"
)

type LoginUserCommand struct {
	Name           string
	DeviceDto      DeviceDto
	ExternalId     string
	Surname        string
	CredentialsDto CredentialsDto
}

type CredentialsDto struct {
	AccessToken              string
	RefreshToken             string
	Scopes                   string
	AccessTokenExpiresAtUtc  time.Time
	RefreshTokenExpiresAtUtc time.Time
	IssuedAt                 time.Time
}
type DeviceDto struct {
	DeviceId   *user.DeviceId
	DeviceType user.DeviceType
}
