package external_credentials

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type ExternalCredentials struct {
	shared.AggregateRoot[user_id.UserId]
	externalId               string
	accessToken              string
	refreshToken             string
	musicProvider            music_provider.MusicProvider
	accessTokenExpiresAtUtc  time.Time
	refreshTokenExpiresAtUtc time.Time
	issuedAtUtc              time.Time
}

func NewExternalCredentials(
	userId user_id.UserId,
	accessToken,
	refreshToken,
	externalId string,
	musicProvider music_provider.MusicProvider,
	accessTokenExpiration,
	refreshTokenExpiration time.Time,
) (*ExternalCredentials, error) {
	if accessTokenExpiration.Before(time.Now().UTC()) {
		return nil, validation_domain_error.NewValidationDomainError(
			"ExternalCredentials.AccessToken.Expired",
			"Access token expiration time must be in the future",
		)
	}
	if refreshTokenExpiration.Before(time.Now().UTC()) {
		return nil, validation_domain_error.NewValidationDomainError(
			"ExternalCredentials.RefreshToken.Expired",
			"Refresh token expiration time must be in the future",
		)
	}
	if accessToken == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"ExternalCredentials.AccessToken",
			"access token")
	}
	if refreshToken == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"ExternalCredentials.RefreshToken",
			"refresh token")
	}
	if externalId == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"ExternalCredentials.ExternalId",
			"external id")
	}
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
	return creds, nil
}

func (cr *ExternalCredentials) AccessToken() string {
	return cr.accessToken
}

func (cr *ExternalCredentials) RefreshToken() string {
	return cr.refreshToken
}

func (cr *ExternalCredentials) MusicProvider() music_provider.MusicProvider {
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

func HydrateExternalCredentials(id user_id.UserId,
	accessToken,
	refreshToken,
	externalId string,
	musicProvider music_provider.MusicProvider,
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
