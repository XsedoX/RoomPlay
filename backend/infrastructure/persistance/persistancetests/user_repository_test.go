package persistancetests

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepositoryAdd(t *testing.T) {
	txx, ctx := integration_tests.GetTxxAndCtx(t, false)

	repo := persistance.NewUserRepository()

	userID := user.Id(uuid.New())
	deviceID := user.DeviceId(uuid.New())
	roomID := shared.RoomId(uuid.New())

	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, uuid.UUID(roomID), time.Now().UTC())
	require.NoError(t, err)

	device := user.HydrateDevice(
		deviceID,
		"My Device",
		user.Mobile,
		false,
		user.Online,
		time.Now().UTC(),
	)

	u := user.HydrateUser(
		userID,
		"John",
		"Doe",
		nil, // Role
		&roomID,
		[]user.Device{*device},
		nil, // BoostUsedAt
	)

	// Act
	err = repo.Add(ctx, u, txx)
	require.NoError(t, err)

	// Assert User
	var userDb daos.UserDao
	err = txx.GetContext(ctx, &userDb, "SELECT * FROM users WHERE id = $1", uuid.UUID(userID))
	require.NoError(t, err)
	assert.Equal(t, "John", userDb.Name)
	assert.Equal(t, "Doe", userDb.Surname)

	// Assert Device
	var deviceDb daos.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", uuid.UUID(deviceID))
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, uuid.UUID(userID), deviceDb.UserId)
}

func TestUserRepositoryGetUserById(t *testing.T) {
	txx, ctx := integration_tests.GetTxxAndCtx(t, false)

	repo := persistance.NewUserRepository()

	// Setup Data
	userID := uuid.New()
	deviceID := uuid.New()

	_, err := txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Jane', 'Doe')`, userID)
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc) VALUES ($1, 'Device 2', false, 'mobile', $2, 'offline', $3)`, deviceID, userID, time.Now().UTC())
	require.NoError(t, err)

	// Act
	u, err := repo.GetUserById(ctx, user.Id(userID), txx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, userID, uuid.UUID(u.Id()))
	assert.Equal(t, "Jane", u.FullName().Name())
	assert.Len(t, u.Devices(), 1)
	assert.Equal(t, deviceID, uuid.UUID(u.Devices()[0].Id()))
}

func TestUserRepositoryUpdate(t *testing.T) {
	txx, ctx := integration_tests.GetTxxAndCtx(t, false)

	repo := persistance.NewUserRepository()

	// Setup Data
	userID := uuid.New()
	roomID := uuid.New()

	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Bob', 'Smith')`, userID)
	require.NoError(t, err)

	// Create domain object to update
	rID := shared.RoomId(roomID)
	role := user.Member
	boostTime := time.Now().UTC()
	deviceID := user.DeviceId(uuid.New())

	device := user.HydrateDevice(
		deviceID,
		"My Device",
		user.Mobile,
		false,
		user.Online,
		time.Now().UTC(),
	)

	u := user.HydrateUser(
		user.Id(userID),
		"Bobby", // Changed name
		"Smith",
		&role, // Added role
		&rID,  // Added room
		[]user.Device{*device},
		&boostTime, // Added boost
	)

	// Act
	err = repo.Update(ctx, u, txx)
	require.NoError(t, err)

	// Assert
	var userDb daos.UserDao
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
	var deviceDb daos.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", deviceID)
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, userID, deviceDb.UserId)
	assert.Equal(t, deviceDb.Type, device.DeviceType().String())
	assert.Equal(t, deviceDb.State, user.Online.String())
}

func TestGetUserByExternalIdSuccess(t *testing.T) {
	txx, ctx := integration_tests.GetTxxAndCtx(t, false)
	repo := persistance.NewUserRepository()
	usersExternalId := integration_tests.SeedData.ExternalCredentials[0].ExternalId()

	user, repoErr := repo.GetUserByExternalId(ctx, usersExternalId, txx)
	require.NoError(t, repoErr)
	require.Equal(t, integration_tests.SeedData.Users[0].Id(), user.Id())
	require.Equal(t, integration_tests.SeedData.Users[0].FullName().Name(), user.FullName().Name())
}
