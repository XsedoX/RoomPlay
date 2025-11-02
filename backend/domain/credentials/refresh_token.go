package credentials

import (
	"time"

	"xsedox.com/main/domain/shared"
)

const RefreshTokenExpirationTime = time.Hour * 24 * 7 // 7 days
const TokenExpiredCode = "Refresh.Token.Expired"

type RefreshToken struct {
	shared.AggregateRoot[shared.UserId]
	deviceId     shared.DeviceId
	refreshToken string
	expiresAtUtc time.Time
	issuedAtUtc  time.Time
}

func NewRefreshToken(userId shared.UserId, deviceId shared.DeviceId, refreshToken string) *RefreshToken {
	rt := &RefreshToken{
		refreshToken: refreshToken,
		deviceId:     deviceId,
		expiresAtUtc: time.Now().Add(RefreshTokenExpirationTime).UTC(),
		issuedAtUtc:  time.Now().UTC(),
	}
	rt.SetId(userId)
	return rt
}
func (r RefreshToken) RefreshToken() string {
	return r.refreshToken
}

func (r RefreshToken) ExpiresAtUtc() time.Time {
	return r.expiresAtUtc
}

func (r RefreshToken) IssuedAtUtc() time.Time {
	return r.issuedAtUtc
}
func (r RefreshToken) DeviceId() shared.DeviceId {
	return r.deviceId
}
func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAtUtc().Sub(time.Now().UTC()) <= 0
}
