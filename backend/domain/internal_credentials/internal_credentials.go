package internal_credentials

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

const (
	RefreshTokenExpirationTime = time.Hour * 24 * 7 // 7 days
)

type InternalCredentials struct {
	userSession  shared.AggregateRoot[user_session.UserSession]
	refreshToken token.Token
	issuedAtUtc  time.Time
}

func NewInternalCredentials(
	userSession user_session.UserSession,
	refreshTokenValue string,
) (*InternalCredentials, error) {
	refreshToken, err := token.NewToken(refreshTokenValue, time.Now().Add(RefreshTokenExpirationTime).UTC())
	if err != nil {
		return nil, err
	}
	rt := &InternalCredentials{
		refreshToken: *refreshToken,
		issuedAtUtc:  time.Now().UTC(),
	}
	rt.userSession.SetId(userSession)
	return rt, nil
}

func (r InternalCredentials) RefreshToken() token.Token {
	return r.refreshToken
}

func (r InternalCredentials) ExpiresAtUtc() time.Time {
	return r.refreshToken.ExpiresAtUtc().UTC()
}

func (r InternalCredentials) IssuedAtUtc() time.Time {
	return r.issuedAtUtc.UTC()
}

func (r InternalCredentials) UserId() user_id.UserId {
	userSession := r.userSession.Id()
	return userSession.UserId()
}

func (r InternalCredentials) DeviceId() device_id.DeviceId {
	userSession := r.userSession.Id()
	return userSession.DeviceId()
}

func (r InternalCredentials) IsExpired() bool {
	return r.ExpiresAtUtc().Sub(time.Now().UTC()) <= 0
}

func (r InternalCredentials) UserSession() user_session.UserSession {
	return r.userSession.Id()
}

func HydrateInternalCredentials(
	userSession user_session.UserSession,
	refreshTokenValue string,
	expiresAtUtc time.Time,
	issuedAtUtc time.Time,
) *InternalCredentials {
	refreshToken := token.HydrateToken(refreshTokenValue, expiresAtUtc)
	result := &InternalCredentials{
		refreshToken: *refreshToken,
		issuedAtUtc:  issuedAtUtc.UTC(),
	}
	result.userSession.SetId(userSession)
	return result
}
