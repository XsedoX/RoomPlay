package register

import (
	"time"

	"xsedox.com/main/domain/device"
)

type UserCommand struct {
	Name           string
	DeviceType     device.Type
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
