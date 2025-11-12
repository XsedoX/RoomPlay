package register_user_command

import (
	"time"

	"xsedox.com/main/domain/user"
)

type RegisterUserCommand struct {
	Name           string
	DeviceType     user.DeviceType
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
