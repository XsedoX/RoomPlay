package refresh_access_token_dto

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type RefreshAccessTokenDto struct {
	UserId                  user_id.UserId `json:"user_id"`
	AccessToken             string         `json:"access_token"`
	AccessTokenExpiresAtUtc time.Time      `json:"access_token_expires_at_utc"`
}
