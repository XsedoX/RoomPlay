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
	AccessToken           string
	RefreshToken          string
	Scopes                string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	IssuedAt              time.Time
}
