package credentials

import (
	"time"

	"xsedox.com/main/domain/shared"
)

const RefreshTokenExpirationTime = time.Hour * 24 * 7 // 7 days

type RefreshToken struct {
	shared.AggregateRoot[shared.UserId]
	deviceId       shared.DeviceId
	refreshToken   []byte
	expirationTime time.Time
	issuedAt       time.Time
}

func NewRefreshToken(userId shared.UserId, deviceId shared.DeviceId, refreshToken []byte) *RefreshToken {
	rt := &RefreshToken{
		refreshToken:   refreshToken,
		deviceId:       deviceId,
		expirationTime: time.Now().Add(RefreshTokenExpirationTime),
		issuedAt:       time.Now(),
	}
	rt.SetId(userId)
	return rt
}
func (r RefreshToken) RefreshToken() []byte {
	return r.refreshToken
}

func (r RefreshToken) ExpirationTime() time.Time {
	return r.expirationTime
}

func (r RefreshToken) IssuedAt() time.Time {
	return r.issuedAt
}
func (r RefreshToken) DeviceId() shared.DeviceId {
	return r.deviceId
}
