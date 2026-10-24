package external_credentials

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type ExternalCredentials struct {
	shared.AggregateRoot[user_id.UserId]
	externalId             string
	accessToken            token.Token
	refreshToken           string
	musicProvider          music_provider.MusicProvider
	accessTokenIssuedAtUtc time.Time
}

func NewExternalCredentials(
	userId user_id.UserId,
	accessTokenString,
	refreshTokenString,
	externalId string,
	musicProvider music_provider.MusicProvider,
	accessTokenExpiration time.Time,
) (*ExternalCredentials, error) {
	accessToken, err := token.NewToken(accessTokenString, accessTokenExpiration)
	if err != nil {
		return nil, err
	}
	if externalId == "" {
		return nil, domain_errors.NewExternalCredentialsExternalIdEmptyError()
	}
	creds := &ExternalCredentials{
		accessToken:            *accessToken,
		refreshToken:           refreshTokenString,
		externalId:             externalId,
		musicProvider:          musicProvider,
		accessTokenIssuedAtUtc: time.Now().UTC(),
	}
	creds.SetId(userId)
	return creds, nil
}

func (cr *ExternalCredentials) UpdateAccessToken(newToken token.Token) {
	cr.accessToken = newToken
}

func (cr *ExternalCredentials) GetAccessToken() token.Token {
	return cr.accessToken
}

func (cr *ExternalCredentials) GetRefreshToken() string {
	return cr.refreshToken
}

func (cr *ExternalCredentials) MusicProvider() music_provider.MusicProvider {
	return cr.musicProvider
}

func (cr *ExternalCredentials) AccessTokenExpiresAtUtc() time.Time {
	return cr.accessToken.ExpiresAtUtc().UTC()
}

func (cr *ExternalCredentials) IssuedAtUtc() time.Time {
	return cr.accessTokenIssuedAtUtc.UTC()
}

func (cr *ExternalCredentials) ExternalId() string {
	return cr.externalId
}

func HydrateExternalCredentials(id user_id.UserId,
	accessToken,
	refreshToken,
	externalId string,
	musicProvider music_provider.MusicProvider,
	accessTokenExpiresAtUtc,
	issuedAtUtc time.Time,
) *ExternalCredentials {
	accessTokenVO := token.HydrateToken(accessToken, accessTokenExpiresAtUtc)

	creds := &ExternalCredentials{
		accessToken:            *accessTokenVO,
		refreshToken:           refreshToken,
		externalId:             externalId,
		musicProvider:          musicProvider,
		accessTokenIssuedAtUtc: issuedAtUtc.UTC(),
	}
	creds.SetId(id)
	return creds
}
