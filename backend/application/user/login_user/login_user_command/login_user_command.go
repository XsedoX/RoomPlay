package login_user_command

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
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
	MusicProvider            music_provider.MusicProvider
	AccessTokenExpiresAtUtc  time.Time
	RefreshTokenExpiresAtUtc time.Time
	IssuedAt                 time.Time
}
type DeviceDto struct {
	DeviceId   *device_id.DeviceId
	DeviceType device_type.DeviceType
}
