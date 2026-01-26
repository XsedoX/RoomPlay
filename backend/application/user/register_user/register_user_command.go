package register_user

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type RegisterUserCommand struct {
	Name           string
	DeviceType     user.DeviceType
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
