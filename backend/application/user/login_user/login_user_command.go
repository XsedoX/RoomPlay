package login_user

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type LoginUserCommand struct {
	Name           string
	DeviceDto      DeviceDto
	Surname        string
	CredentialsDto CredentialsDto
}

type CredentialsDto struct {
	ExternalId               string
	AccessToken              string
	RefreshToken             string
	MusicProvider            credentials.MusicProvider
	AccessTokenExpiresAtUtc  time.Time
	RefreshTokenExpiresAtUtc time.Time
	IssuedAt                 time.Time
}
type DeviceDto struct {
	DeviceId   *user.DeviceId
	DeviceType user.DeviceType
}
