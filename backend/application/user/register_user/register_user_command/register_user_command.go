package register_user_command

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
)

type RegisterUserCommand struct {
	Name           string
	DeviceType     device_type.DeviceType
	Surname        string
	CredentialsDto CredentialsDto
}

type CredentialsDto struct {
	ExternalId               string
	AccessToken              string
	RefreshToken             string
	MusicProvider            music_provider.MusicProvider
	AccessTokenExpiresAtUtc  time.Time
	RefreshTokenExpiresAtUtc time.Time
	IssuedAt                 time.Time
}
