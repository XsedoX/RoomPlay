package external_credentials

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewExternalCredentialsSuccess(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := faker.Jwt()
	accessTokenExpiration := time.Now().Add(24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour).UTC()
	accessToken := faker.Jwt()
	externalId := faker.UUIDDigit()

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.NoError(t, error)
	require.Equal(t, userId, externalCred.Id())
	require.Equal(t, externalId, externalCred.ExternalId())
	require.Equal(t, music_provider.YouTube, externalCred.MusicProvider())
	require.Equal(t, accessToken, externalCred.GetAccessToken().Value())
	require.Equal(t, refreshToken, externalCred.GetRefreshToken().Value())
	require.Equal(t, accessTokenExpiration, externalCred.AccessTokenExpiresAtUtc())
	require.Equal(t, refreshTokenExpiration, externalCred.RefreshTokenExpiresAtUtc())
	require.WithinDuration(t, time.Now().UTC(), externalCred.IssuedAtUtc(), time.Second)
}

func TestNewExternalCredentialsAccessTokenExpired(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := faker.Jwt()
	accessTokenExpiration := time.Now().Add(-24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour).UTC()
	accessToken := faker.Jwt()
	externalId := faker.UUIDDigit()

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Nil(t, externalCred)
	require.Error(t, error)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "Token.Expired", castedError.Code)
	require.Equal(t, "Token expiration time must be in the future", castedError.Description)
}

func TestNewExternalCredentialsRefreshTokenExpired(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := faker.Jwt()
	accessTokenExpiration := time.Now().Add(24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(-30 * 24 * time.Hour).UTC()
	accessToken := faker.Jwt()
	externalId := faker.UUIDDigit()

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Nil(t, externalCred)
	require.Error(t, error)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "Token.Expired", castedError.Code)
	require.Equal(t, "Token expiration time must be in the future", castedError.Description)
}

func TestNewExternalCredentialsAccessTokenEmpty(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := faker.Jwt()
	accessTokenExpiration := time.Now().Add(24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour).UTC()
	accessToken := ""
	externalId := faker.UUIDDigit()

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Nil(t, externalCred)
	require.Error(t, error)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "Token.Empty", castedError.Code)
	require.Equal(t, "Token cannot be empty", castedError.Description)
}

func TestNewExternalCredentialsRefreshTokenEmpty(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := ""
	accessTokenExpiration := time.Now().Add(24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour).UTC()
	accessToken := faker.Jwt()
	externalId := faker.UUIDDigit()

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Nil(t, externalCred)
	require.Error(t, error)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "Token.Empty", castedError.Code)
	require.Equal(t, "Token cannot be empty", castedError.Description)
}

func TestNewExternalCredentialsExternalIdEmpty(t *testing.T) {
	userId := user_id.NewUserId()
	refreshToken := faker.Jwt()
	accessTokenExpiration := time.Now().Add(24 * time.Hour).UTC()
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour).UTC()
	accessToken := faker.Jwt()
	externalId := ""

	externalCred, error := NewExternalCredentials(userId,
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiration,
		refreshTokenExpiration,
	)

	require.Nil(t, externalCred)
	require.Error(t, error)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.ExternalId.Empty", castedError.Code)
	require.Equal(t, "External id cannot be empty", castedError.Description)
}
