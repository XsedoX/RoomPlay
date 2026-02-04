package external_credentials_repository

import (
	"context"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeDatabaseContainer()
	tests_initializer.RunTestsWithDatabase(m)
}

func setupMocks(t *testing.T) (*sqlx.Tx,
	context.Context,
	*mock_encrypter.MockEncrypter,
) {
	txx, ctx := tests_initializer.GetTxxAndCtx(t, false)
	mockEncrypter := new(mock_encrypter.MockEncrypter)

	defer mockEncrypter.AssertExpectations(t)

	return txx, ctx, mockEncrypter
}

func TestExternalCredentialsRepositoryGrant(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := NewExternalCredentialsRepository(mockEncrypter)
	// Get a user from the seeded database
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err, "failed to find a user in the database")

	accessToken := "access_token_123"
	refreshToken := "refresh_token_123"
	accessTokenExpiresAt := time.Now().Add(1 * time.Hour).UTC()
	refreshTokenExpiresAt := time.Now().Add(24 * time.Hour).UTC()
	externalId := faker.UUIDDigit()

	creds, extCredsErr := external_credentials.NewExternalCredentials(
		user_id.UserId(userID),
		accessToken,
		refreshToken,
		externalId,
		music_provider.YouTube,
		accessTokenExpiresAt,
		refreshTokenExpiresAt,
	)
	require.NoError(t, extCredsErr, "failed to create ExternalCredentials domain object")

	// Setup mock expectations
	mockEncrypter.On("Encrypt", accessToken).Return([]byte("encrypted_"+accessToken), nil)
	mockEncrypter.On("Encrypt", refreshToken).Return([]byte("encrypted_"+refreshToken), nil)

	// Act
	err = repo.Grant(ctx, creds, txx)
	require.NoError(t, err)

	// Assert
	var storedCreds struct {
		UserID                   uuid.UUID `db:"user_id"`
		AccessToken              []byte    `db:"access_token"`
		RefreshToken             []byte    `db:"refresh_token"`
		ExternalId               string    `db:"external_id"`
		AccessTokenExpiresAtUtc  time.Time `db:"access_token_expires_at_utc"`
		RefreshTokenExpiresAtUtc time.Time `db:"refresh_token_expires_at_utc"`
		IssuedAtUtc              time.Time `db:"issued_at_utc"`
		MusicProvider            string    `db:"music_provider"`
	}

	err = txx.GetContext(ctx, &storedCreds, "SELECT * FROM users_external_credentials WHERE user_id = $1", userID)
	require.NoError(t, err)

	assert.Equal(t, userID, storedCreds.UserID)
	assert.Equal(t, externalId, storedCreds.ExternalId)
	assert.Equal(t, []byte("encrypted_"+accessToken), storedCreds.AccessToken)
	assert.Equal(t, []byte("encrypted_"+refreshToken), storedCreds.RefreshToken)
	assert.Equal(t, music_provider.YouTube.String(), storedCreds.MusicProvider)
	// Compare times with small delta to account for DB roundtrip precision
	assert.WithinDuration(t, accessTokenExpiresAt, storedCreds.AccessTokenExpiresAtUtc, time.Second)
	assert.WithinDuration(t, refreshTokenExpiresAt, storedCreds.RefreshTokenExpiresAtUtc, time.Second)
	assert.WithinDuration(t, time.Now().UTC(), storedCreds.IssuedAtUtc, 5*time.Second)
}
