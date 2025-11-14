package credentials

import (
	"time"

	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

const (
	RefreshTokenExpirationTime = time.Hour * 24 * 7 // 7 days
)

type RefreshToken struct {
	shared.AggregateRoot[user.Id]
	deviceId     user.DeviceId
	refreshToken string
	expiresAtUtc time.Time
	issuedAtUtc  time.Time
}

func NewRefreshToken(userId user.Id, deviceId user.DeviceId, refreshToken string) *RefreshToken {
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
func (r RefreshToken) DeviceId() user.DeviceId {
	return r.deviceId
}
func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAtUtc().Sub(time.Now().UTC()) <= 0
}
func HydrateRefreshToken(userId user.Id,
	deviceId user.DeviceId,
	refreshToken string,
	expiresAtUtc time.Time,
	issuedAtUtc time.Time) *RefreshToken {
	result := &RefreshToken{
		deviceId:     deviceId,
		refreshToken: refreshToken,
		expiresAtUtc: expiresAtUtc.UTC(),
		issuedAtUtc:  issuedAtUtc.UTC(),
	}
	result.SetId(userId)
	return result
}
