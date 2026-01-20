package credentials

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewExternalCredentials(t *testing.T) {
	userId := user.Id(uuid.New())
	refreshToken := faker.Jwt()
	scopes := "read offline"
	accessTokenExpiration := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
	refreshTokenExpiration := time.Date(2040, 12, 31, 23, 59, 59, 0, time.UTC)
	accessToken := faker.Jwt()

	externalCred := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		scopes,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Equal(t, userId, externalCred.Id())
	require.Equal(t, accessToken, externalCred.AccessToken())
	require.Equal(t, refreshToken, externalCred.RefreshToken())
	require.Equal(t, []string{"read", "offline"}, externalCred.Scopes())
	require.Equal(t, accessTokenExpiration, externalCred.AccessTokenExpiresAtUtc())
	require.Equal(t, refreshTokenExpiration, externalCred.RefreshTokenExpiresAtUtc())
	require.WithinDuration(t, time.Now().UTC(), externalCred.IssuedAtUtc(), time.Second)
}
