package persistance

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
	infrastructuretest "xsedox.com/main/test_helpers/infrustructure_test"
	"xsedox.com/main/test_helpers/infrustructure_test/authentication_mocks"
)

var (
	pgContainer *infrastructuretest.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgContainer, err = infrastructuretest.SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres: %v", err)
	}

	projectRoot, err := infrastructuretest.FindProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
	if err := pgContainer.ApplySchema(ctx, schemaPath); err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}

	// Seed database once
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	seeder := infrastructuretest.NewSeeder(dbx)
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	code := m.Run()

	if err := pgContainer.Teardown(ctx); err != nil {
		log.Printf("failed to teardown postgres: %v", err)
	}

	os.Exit(code)
}

func TestExternalCredentialsRepository_Grant(t *testing.T) {
	ctx := context.Background()

	// Start transaction
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := NewExternalCredentialsRepository(mockEncrypter)

	// Get a user from the seeded database
	var userID uuid.UUID
	err = txx.QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&userID)
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
