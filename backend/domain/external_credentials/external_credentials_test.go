package external_credentials

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
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
	require.Equal(t, accessToken, externalCred.AccessToken())
	require.Equal(t, refreshToken, externalCred.RefreshToken())
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
	castedError, ok := error.(*validation_domain_error.ValidationDomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.AccessToken.Expired", castedError.Code)
	require.Equal(t, "Access token expiration time must be in the future", castedError.Description)
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
	castedError, ok := error.(*validation_domain_error.ValidationDomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.RefreshToken.Expired", castedError.Code)
	require.Equal(t, "Refresh token expiration time must be in the future", castedError.Description)
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
	castedError, ok := error.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.AccessToken.EmptyString", castedError.Code)
	require.Equal(t, "The field 'access token' cannot be an empty string.", castedError.Description)
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
	castedError, ok := error.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.RefreshToken.EmptyString", castedError.Code)
	require.Equal(t, "The field 'refresh token' cannot be an empty string.", castedError.Description)
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
	castedError, ok := error.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "ExternalCredentials.ExternalId.EmptyString", castedError.Code)
	require.Equal(t, "The field 'external id' cannot be an empty string.", castedError.Description)
}
