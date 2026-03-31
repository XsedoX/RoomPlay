package internal_credentials_repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/go-faker/faker/v4"
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
	*internal_credentials_repository.InternalCredentialsRepository,
) {
	txx, ctx := tests_initializer.GetTxxAndCtx(t, false)
	mockEncrypter := new(mock_encrypter.MockEncrypter)
	repo := internal_credentials_repository.NewInternalCredentialsRepository(mockEncrypter)

	defer mockEncrypter.AssertExpectations(t)

	return txx, ctx, mockEncrypter, repo
}

func TestInternalCredentialsRepositoryAssignNewToken(t *testing.T) {
	txx,
		ctx,
		mockEncrypter,
		repo := setupMocks(t)

	tokenValue := faker.UUIDDigit()
	encryptedToken := []byte("encrypted_" + tokenValue)
	expiresAt := time.Now().Add(24 * time.Hour).UTC()
	issuedAt := time.Now().UTC()

	mockEncrypter.On("Hash", tokenValue).Return(encryptedToken)

	notLoggedInUser := seeder.SeedData.Users[1]
	userSession := user_session.NewUserSession(
		notLoggedInUser.Id(),
		notLoggedInUser.Devices()[0].Id(),
	)
	internalCredentials := internal_credentials.HydrateInternalCredentials(
		*userSession,
		tokenValue,
		expiresAt,
		issuedAt,
	)

	repoErr := repo.AssignNewToken(ctx, internalCredentials, txx)
	require.NoError(t, repoErr)
	var tokenFromDb internal_credentials_dao.InternalCredentialsDao
	err := txx.GetContext(ctx, &tokenFromDb, `
		SELECT *
		FROM users_internal_credentials
		WHERE user_id = $1 AND device_id = $2
		`, userSession.UserId().ToUuid(),
		userSession.DeviceId().ToUuid(),
	)
	require.NoError(t, err)
	require.Equal(t, tokenFromDb.RefreshToken, encryptedToken)
	mockEncrypter.AssertNumberOfCalls(t, "Hash", 1)
}

func TestInternalCredentialsRepositoryGetTokenByValue(t *testing.T) {
	txx,
		ctx,
		mockEncrypter,
		repo := setupMocks(t)

	configuration := mock_configuration.MockConfiguration{}
	realEncrypter := encryper.NewEncrypter(&configuration)
	existingToken := seeder.SeedData.InternalCredentials[0].RefreshToken()
	encryptedExistingToken := realEncrypter.Hash(existingToken)
	mockEncrypter.On("Hash", existingToken).Return(encryptedExistingToken)

	internalCredentials, err := repo.GetTokenByValue(ctx, existingToken, txx)
	require.NoError(t, err)

	require.Equal(t, encryptedExistingToken, []byte(internalCredentials.RefreshToken()))
	require.WithinDuration(t, seeder.SeedData.InternalCredentials[0].ExpiresAtUtc(), internalCredentials.ExpiresAtUtc(), time.Second)
	require.WithinDuration(t, seeder.SeedData.InternalCredentials[0].IssuedAtUtc(), internalCredentials.IssuedAtUtc(), time.Second)
	require.Equal(t, seeder.SeedData.InternalCredentials[0].UserId(), internalCredentials.UserId())
	require.Equal(t, seeder.SeedData.InternalCredentials[0].DeviceId(), internalCredentials.DeviceId())

	mockEncrypter.AssertNumberOfCalls(t, "Hash", 1)
}

func TestInternalCredentialsRepositoryRetireTokenByUserSession(t *testing.T) {
	txx,
		ctx,
		_,
		repo := setupMocks(t)

	// Act
	userSession := seeder.SeedData.InternalCredentials[0].UserSession()
	err := repo.RetireTokenByUserSession(ctx, userSession, txx)
	require.NoError(t, err)

	// Assert
	var count int
	err = txx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM users_internal_credentials
		WHERE user_id = $1 AND device_id = $2
	`, userSession.UserId().ToUuid(),
		userSession.DeviceId().ToUuid(),
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInternalCredentialsRepositoryRetireAllTokensByUserId(t *testing.T) {
	txx,
		ctx,
		_,
		repo := setupMocks(t)

	// Act
	uID := seeder.SeedData.InternalCredentials[0].UserId()
	err := repo.RetireAllTokensByUserId(ctx, &uID, txx)
	require.NoError(t, err)

	// Assert
	var count int
	err = txx.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM users_internal_credentials
		WHERE user_id = $1
	`, uID.ToUuid(),
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
