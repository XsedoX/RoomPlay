package credentials

import (
	"strings"
	"time"

	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type External struct {
	shared.AggregateRoot[user.Id]
	accessToken              string
	refreshToken             string
	scopes                   []string
	accessTokenExpiresAtUtc  time.Time
	refreshTokenExpiresAtUtc time.Time
	issuedAtUtc              time.Time
}

func NewExternalCredentials(userId user.Id, accessToken, refreshToken, scopes string, accessTokenExpiration, refreshTokenExpiration time.Time) *External {
	creds := &External{
		accessToken:              accessToken,
		refreshToken:             refreshToken,
		scopes:                   strings.Split(scopes, " "),
		accessTokenExpiresAtUtc:  accessTokenExpiration,
		refreshTokenExpiresAtUtc: refreshTokenExpiration,
		issuedAtUtc:              time.Now(),
	}
	creds.SetId(userId)
	return creds
}
func (cr *External) AccessToken() string {
	return cr.accessToken
}
func (cr *External) RefreshToken() string {
	return cr.refreshToken
}
func (cr *External) Scopes() []string {
	return cr.scopes
}
func (cr *External) AccessTokenExpiresAtUtc() time.Time {
	return cr.accessTokenExpiresAtUtc
}
func (cr *External) RefreshTokenExpiresAtUtc() time.Time {
	return cr.refreshTokenExpiresAtUtc
}
func (cr *External) IssuedAtUtc() time.Time {
	return cr.issuedAtUtc
}
