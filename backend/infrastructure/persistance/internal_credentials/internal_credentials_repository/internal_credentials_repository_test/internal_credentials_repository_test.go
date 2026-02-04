package internal_credentials_repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
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

func TestInternalCredentialsRepositoryAssignNewToken(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := internal_credentials_repository.NewInternalCredentialsRepository(mockEncrypter)
	// Get a user from the seeded database
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err, "failed to find a user in the database")

	// Create a device for the user (if not exists, but seeder should have created one)
	// Let's just insert a new device to be sure and clean
	deviceID := uuid.New()
	_, err = txx.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, deviceID, "Test Device 2", false, "mobile", userID, "online", time.Now().UTC())
	require.NoError(t, err)

	tokenValue := "refresh_token_abc"
	expiresAt := time.Now().Add(24 * time.Hour).UTC()
	issuedAt := time.Now().UTC()

	userSession := user_session.NewUserSession(
		user_id.UserId(userID),
		device_id.DeviceId(deviceID),
	)
	internalCredentials := internal_credentials.HydrateInternalCredentials(
		*userSession,
		tokenValue,
		expiresAt,
		issuedAt,
	)

	// Mock expectations
	encryptedToken := []byte("encrypted_" + tokenValue)
	mockEncrypter.On("Hash", tokenValue).Return(encryptedToken)

	// Act
	err = repo.AssignNewToken(ctx, internalCredentials, txx)
	require.NoError(t, err)

	// Assert
	var storedToken struct {
		UserID       uuid.UUID `db:"user_id"`
		DeviceId     uuid.UUID `db:"device_id"`
		RefreshToken []byte    `db:"refresh_token"`
		ExpiresAtUtc time.Time `db:"expires_at_utc"`
		IssuedAtUtc  time.Time `db:"issued_at_utc"`
	}

	err = txx.GetContext(ctx, &storedToken, "SELECT * FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2", userID, deviceID)
	require.NoError(t, err)

	assert.Equal(t, userID, storedToken.UserID)
	assert.Equal(t, deviceID, storedToken.DeviceId)
	assert.Equal(t, encryptedToken, storedToken.RefreshToken)
	assert.WithinDuration(t, expiresAt, storedToken.ExpiresAtUtc, time.Second)
	assert.WithinDuration(t, issuedAt, storedToken.IssuedAtUtc, time.Second)
}

func TestInternalCredentialsRepositoryGetTokenByValue(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := internal_credentials_repository.NewInternalCredentialsRepository(mockEncrypter)

	// Get a user from the seeded database
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err, "failed to find a user in the database")

	// Create device
	deviceID := uuid.New()
	_, err = txx.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, deviceID, "Test Device Get", false, "mobile", userID, "online", time.Now().UTC())
	require.NoError(t, err)

	tokenValue := "refresh_token_xyz"
	encryptedToken := []byte("encrypted_" + tokenValue)
	expiresAt := time.Now().Add(24 * time.Hour).UTC()
	issuedAt := time.Now().UTC()

	// Insert token directly into DB
	_, err = txx.ExecContext(ctx, `
		INSERT INTO users_refresh_tokens (user_id, device_id, refresh_token, expires_at_utc, issued_at_utc)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, deviceID, encryptedToken, expiresAt, issuedAt)
	require.NoError(t, err)

	// Mock expectations
	mockEncrypter.On("Hash", tokenValue).Return(encryptedToken)

	// Act
	retrievedToken, err := repo.GetTokenByValue(ctx, tokenValue, txx)
	require.NoError(t, err)

	// Assert
	assert.NotNil(t, retrievedToken)
	assert.Equal(t, user_id.UserId(userID), retrievedToken.Id())
	assert.Equal(t, device_id.DeviceId(deviceID), retrievedToken.DeviceId())
	assert.Equal(t, string(encryptedToken), retrievedToken.RefreshToken())
}

func TestInternalCredentialsRepositoryRetireTokenByUserIdAndDeviceId(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)

	repo := internal_credentials_repository.NewInternalCredentialsRepository(mockEncrypter)

	// Get a user
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err)

	// Create device
	deviceID := uuid.New()
	_, err = txx.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, deviceID, "Test Device Retire", false, "mobile", userID, "online", time.Now().UTC())
	require.NoError(t, err)

	// Insert token
	_, err = txx.ExecContext(ctx, `
		INSERT INTO users_refresh_tokens (user_id, device_id, refresh_token, expires_at_utc, issued_at_utc)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, deviceID, []byte("some_token"), time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	// Act
	uID := user_id.UserId(userID)
	dID := device_id.DeviceId(deviceID)
	err = repo.RetireTokenByUserIdAndDeviceId(ctx, &uID, &dID, txx)
	require.NoError(t, err)

	// Assert
	var count int
	err = txx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2", userID, deviceID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInternalCredentialsRepositoryRetireAllTokensByUserId(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := internal_credentials_repository.NewInternalCredentialsRepository(mockEncrypter)

	// Get a user
	var userID uuid.UUID
	err := txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
	require.NoError(t, err)

	// Create 2 devices
	deviceID1 := uuid.New()
	deviceID2 := uuid.New()
	_, err = txx.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, 'D1', false, 'mobile', $3, 'online', $4), ($2, 'D2', false, 'mobile', $3, 'online', $4)
	`, deviceID1, deviceID2, userID, time.Now().UTC())
	require.NoError(t, err)

	// Insert tokens
	_, err = txx.ExecContext(ctx, `
		INSERT INTO users_refresh_tokens (user_id, device_id, refresh_token, expires_at_utc, issued_at_utc)
		VALUES ($1, $2, $4, $6, $6), ($1, $3, $5, $6, $6)
	`, userID, deviceID1, deviceID2, []byte("token1"), []byte("token2"), time.Now().UTC())
	require.NoError(t, err)

	// Act
	uID := user_id.UserId(userID)
	err = repo.RetireAllTokensByUserId(ctx, &uID, txx)
	require.NoError(t, err)

	// Assert
	var count int
	err = txx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users_refresh_tokens WHERE user_id = $1", userID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
