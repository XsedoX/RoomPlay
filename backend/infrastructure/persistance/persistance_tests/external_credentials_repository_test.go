package persistance_tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/test_helpers/integration_tests"
	"xsedox.com/main/test_helpers/integration_tests/authentication_mocks"
)

func TestExternalCredentialsRepositoryGrant(t *testing.T) {
	txx, ctx := integration_tests.GetTxxAndCtx(t)

	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := persistance.NewExternalCredentialsRepository(mockEncrypter)

	// Get a user from the seeded database
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err, "failed to find a user in the database")

	accessToken := "access_token_123"
	refreshToken := "refresh_token_123"
	accessTokenExpiresAt := time.Now().Add(1 * time.Hour).UTC()
	refreshTokenExpiresAt := time.Now().Add(24 * time.Hour).UTC()

	creds := credentials.NewExternalCredentials(
		user.Id(userID),
		accessToken,
		refreshToken,
		"scope1 scope2",
		accessTokenExpiresAt,
		refreshTokenExpiresAt,
	)

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
		Scope                    string    `db:"scope"`
		AccessTokenExpiresAtUtc  time.Time `db:"access_token_expires_at_utc"`
		RefreshTokenExpiresAtUtc time.Time `db:"refresh_token_expires_at_utc"`
		IssuedAtUtc              time.Time `db:"issued_at_utc"`
	}

	err = txx.GetContext(ctx, &storedCreds, "SELECT * FROM users_external_credentials WHERE user_id = $1", userID)
	require.NoError(t, err)

	assert.Equal(t, userID, storedCreds.UserID)
	assert.Equal(t, []byte("encrypted_"+accessToken), storedCreds.AccessToken)
	assert.Equal(t, []byte("encrypted_"+refreshToken), storedCreds.RefreshToken)
	assert.Equal(t, "scope1 scope2", storedCreds.Scope)
	// Compare times with small delta to account for DB roundtrip precision
	assert.WithinDuration(t, accessTokenExpiresAt, storedCreds.AccessTokenExpiresAtUtc, time.Second)
	assert.WithinDuration(t, refreshTokenExpiresAt, storedCreds.RefreshTokenExpiresAtUtc, time.Second)
	assert.WithinDuration(t, time.Now().UTC(), storedCreds.IssuedAtUtc, 5*time.Second)

	mockEncrypter.AssertExpectations(t)
}
