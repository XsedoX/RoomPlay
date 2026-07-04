package user_repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
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

func TestUserRepositoryAdd(t *testing.T) {
	txx,
		ctx,
		_ := setupMocks(t)

	repo := user_repository.NewUserRepository()

	userID := user_id.NewUserId()
	deviceID := device_id.NewDeviceId()
	roomID := room_id.NewRoomId()

	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, uuid.UUID(roomID), time.Now().UTC())
	require.NoError(t, err)

	device1 := device.HydrateDevice(
		deviceID,
		"My Device",
		device_type.Mobile,
		false,
		device_state.Online,
		time.Now().UTC(),
	)

	u := user.HydrateUser(
		userID,
		"John",
		"Doe",
		nil, // Role
		&roomID,
		[]device.Device{*device1},
		nil, // BoostUsedAt
	)

	// Act
	err = repo.Add(ctx, u, txx)
	require.NoError(t, err)

	// Assert User
	var userDb user_dao.UserDao
	err = txx.GetContext(ctx, &userDb, "SELECT * FROM users WHERE id = $1", uuid.UUID(userID))
	require.NoError(t, err)
	assert.Equal(t, "John", userDb.Name)
	assert.Equal(t, "Doe", userDb.Surname)

	// Assert Device
	var deviceDb device_dao.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", uuid.UUID(deviceID))
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, uuid.UUID(userID), deviceDb.UserId)
}

func TestUserRepositoryGetUserById(t *testing.T) {
	txx,
		ctx,
		_ := setupMocks(t)

	repo := user_repository.NewUserRepository()

	userToTest := seeder.SeedData.Users[1]
	userID := userToTest.Id()
	deviceID := userToTest.Devices()[0].Id()
	// Act
	u, err := repo.GetUserById(ctx, userID, txx)
	require.NoError(t, err)

	// Assert
	obtainedUserId := u.Id()
	assert.Equal(t, userID.ToUuid(), obtainedUserId.ToUuid())
	assert.Equal(t, userToTest.FullName().Name(), u.FullName().Name())
	assert.Len(t, u.Devices(), 1)
	assert.Equal(t, deviceID.ToUuid(), u.Devices()[0].Id().ToUuid())
}

func TestUserRepositoryUpdate(t *testing.T) {
	txx,
		ctx,
		_ := setupMocks(t)
	repo := user_repository.NewUserRepository()

	// Setup Data
	userID := uuid.New()
	roomID := uuid.New()

	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Bob', 'Smith')`, userID)
	require.NoError(t, err)

	// Create domain object to update
	rID := room_id.RoomId(roomID)
	role := user_role.Member
	boostTime := time.Now().UTC()
	deviceID := device_id.NewDeviceId()

	device1 := device.HydrateDevice(
		deviceID,
		"My Device",
		device_type.Mobile,
		false,
		device_state.Online,
		time.Now().UTC(),
	)

	u := user.HydrateUser(
		user_id.UserId(userID),
		"Bobby", // Changed name
		"Smith",
		&role, // Added role
		&rID,  // Added room
		[]device.Device{*device1},
		&boostTime, // Added boost
	)

	// Act
	err = repo.Update(ctx, u, txx)
	require.NoError(t, err)

	// Assert
	var userDb user_dao.UserDao
	err = txx.GetContext(ctx, &userDb, "SELECT * FROM users WHERE id = $1", userID)
	require.NoError(t, err)
	assert.Equal(t, "Bobby", userDb.Name)

	// Check Role
	var roleDb string
	err = txx.QueryRowContext(ctx, "SELECT role FROM users_room_data WHERE user_id = $1", userID).Scan(&roleDb)
	require.NoError(t, err)
	assert.Equal(t, "member", roleDb)

	// Check Boost
	var boostDb time.Time
	err = txx.QueryRowContext(ctx, "SELECT boost_used_at_utc FROM users_room_data WHERE user_id = $1", userID).Scan(&boostDb)
	require.NoError(t, err)
	assert.WithinDuration(t, boostTime, boostDb, time.Second)

	// check device
	var deviceDb device_dao.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", deviceID)
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, userID, deviceDb.UserId)
	assert.Equal(t, deviceDb.Type, device1.DeviceType().String())
	assert.Equal(t, deviceDb.State, device_state.Online.String())
}

func TestGetUserByExternalIdSuccess(t *testing.T) {
	txx,
		ctx,
		_ := setupMocks(t)
	repo := user_repository.NewUserRepository()
	usersExternalId := seeder.SeedData.ExternalCredentials[0].ExternalId()

	user, repoErr := repo.GetUserByExternalId(ctx, usersExternalId, txx)
	require.NoError(t, repoErr)
	require.Equal(t, seeder.SeedData.Users[0].Id(), user.Id())
	require.Equal(t, seeder.SeedData.Users[0].FullName().Name(), user.FullName().Name())
}

func TestCheckIfUserExistByExternalId(t *testing.T) {
	txx,
		ctx,
		_ := setupMocks(t)

	repo := user_repository.NewUserRepository()

	// NOTE: User exists check
	usersExternalId := seeder.SeedData.ExternalCredentials[0].ExternalId()
	userExists := repo.CheckIfUserExistByExternalId(ctx, usersExternalId, txx)
	require.True(t, userExists)

	// NOTE: User does not exist check
	usersExternalIdNonExistent := "this does not exist for sure"
	userExistsNot := repo.CheckIfUserExistByExternalId(ctx, usersExternalIdNonExistent, txx)
	require.False(t, userExistsNot)
}
