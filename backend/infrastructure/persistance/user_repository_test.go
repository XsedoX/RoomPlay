package persistance

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/infrastructure/persistance/daos"
)

func TestUserRepository_Add(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	repo := NewUserRepository()

	userID := user.Id(uuid.New())
	deviceID := user.DeviceId(uuid.New())
	roomID := shared.RoomId(uuid.New())

	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, uuid.UUID(roomID), time.Now().UTC())
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
		"ext-id-1",
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
	assert.Equal(t, "ext-id-1", userDb.ExternalId)

	// Assert Device
	var deviceDb daos.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", uuid.UUID(deviceID))
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, uuid.UUID(userID), deviceDb.UserId)
}

func TestUserRepository_GetUserById(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	repo := NewUserRepository()

	// Setup Data
	userID := uuid.New()
	deviceID := uuid.New()

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-2', 'Jane', 'Doe')`, userID)
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

func TestUserRepository_Update(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	repo := NewUserRepository()

	// Setup Data
	userID := uuid.New()
	roomID := uuid.New()

	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-3', 'Bob', 'Smith')`, userID)
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
		"ext-3",
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

	//check device
	var deviceDb daos.DeviceDao
	err = txx.GetContext(ctx, &deviceDb, "SELECT * FROM devices WHERE id = $1", deviceID)
	require.NoError(t, err)
	assert.Equal(t, "My Device", deviceDb.FriendlyName)
	assert.Equal(t, userID, deviceDb.UserId)
	assert.Equal(t, deviceDb.Type, device.DeviceType().String())
	assert.Equal(t, deviceDb.State, user.Online.String())
}

func TestUserRepository_LeaveRoom(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	repo := NewUserRepository()

	// Setup Data
	userID := uuid.New()
	roomID := uuid.New()

	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-4', 'Alice', 'Wonder')`, userID)
	require.NoError(t, err)

	// Act
	err = repo.LeaveRoom(ctx, user.Id(userID), txx)
	require.NoError(t, err)

	// Assert
	var userDb daos.UserDao
	err = txx.GetContext(ctx, &userDb, "SELECT * FROM users WHERE id = $1", userID)
	require.NoError(t, err)
	assert.Nil(t, userDb.RoomId)
}
