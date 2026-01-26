package credentials

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type ExternalCredentials struct {
	shared.AggregateRoot[user.Id]
	externalId               string
	accessToken              string
	refreshToken             string
	musicProvider            MusicProvider
	accessTokenExpiresAtUtc  time.Time
	refreshTokenExpiresAtUtc time.Time
	issuedAtUtc              time.Time
}

func NewExternalCredentials(userId user.Id,
	accessToken,
	refreshToken,
	externalId string,
	musicProvider MusicProvider,
	accessTokenExpiration,
	refreshTokenExpiration time.Time,
) *ExternalCredentials {
	creds := &ExternalCredentials{
		accessToken:              accessToken,
		refreshToken:             refreshToken,
		externalId:               externalId,
		musicProvider:            musicProvider,
		accessTokenExpiresAtUtc:  accessTokenExpiration,
		refreshTokenExpiresAtUtc: refreshTokenExpiration,
		issuedAtUtc:              time.Now().UTC(),
	}
	creds.SetId(userId)
	return creds
}

func (cr *ExternalCredentials) AccessToken() string {
	return cr.accessToken
}

func (cr *ExternalCredentials) RefreshToken() string {
	return cr.refreshToken
}

func (cr *ExternalCredentials) MusicProvider() MusicProvider {
	return cr.musicProvider
}

func (cr *ExternalCredentials) AccessTokenExpiresAtUtc() time.Time {
	return cr.accessTokenExpiresAtUtc
}

func (cr *ExternalCredentials) RefreshTokenExpiresAtUtc() time.Time {
	return cr.refreshTokenExpiresAtUtc
}

func (cr *ExternalCredentials) IssuedAtUtc() time.Time {
	return cr.issuedAtUtc
}

func (cr *ExternalCredentials) ExternalId() string {
	return cr.externalId
}

func HydrateExternalCredentials(id user.Id,
	accessToken,
	refreshToken,
	externalId string,
	musicProvider MusicProvider,
	accessTokenExpiresAtUtc,
	refreshTokenExpiresAtUtc,
	issuedAtUtc time.Time,
) *ExternalCredentials {
	creds := &ExternalCredentials{
		accessToken:              accessToken,
		refreshToken:             refreshToken,
		externalId:               externalId,
		musicProvider:            musicProvider,
		accessTokenExpiresAtUtc:  accessTokenExpiresAtUtc,
		refreshTokenExpiresAtUtc: refreshTokenExpiresAtUtc,
		issuedAtUtc:              issuedAtUtc,
	}
	creds.SetId(id)
	return creds
}
