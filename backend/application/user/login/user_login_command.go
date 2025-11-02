package login

import (
	"time"

	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/shared"
)

type UserCommand struct {
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
	DeviceId   *shared.DeviceId
	DeviceType device.Type
}
