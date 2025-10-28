package credentials

import (
	"strings"
	"time"

	"xsedox.com/main/domain/shared"
)

type External struct {
	shared.AggregateRoot[shared.UserId]
	accessToken                string
	refreshToken               string
	scopes                     []string
	accessTokenExpirationTime  time.Time
	refreshTokenExpirationTime time.Time
	issuedAt                   time.Time
}

func NewExternalCredentials(userId shared.UserId, accessToken, refreshToken, scopes string, accessTokenExpiration, refreshTokenExpiration time.Time) *External {
	creds := &External{
		accessToken:                accessToken,
		refreshToken:               refreshToken,
		scopes:                     strings.Split(scopes, " "),
		accessTokenExpirationTime:  accessTokenExpiration,
		refreshTokenExpirationTime: refreshTokenExpiration,
		issuedAt:                   time.Now(),
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
func (cr *External) AccessTokenExpirationTime() time.Time {
	return cr.accessTokenExpirationTime
}
func (cr *External) RefreshTokenExpirationTime() time.Time {
	return cr.refreshTokenExpirationTime
}
func (cr *External) IssuedAt() time.Time {
	return cr.issuedAt
}
