package room_repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/slice_extensions"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/room/room_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_dao"
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

func TestRoomRepositoryCreateRoom(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)

	repo := room_repository.NewRoomRepository(mockEncrypter)

	// Setup User
	userID := user_id.NewUserId()
	_, err := txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Host', 'User')`, uuid.UUID(userID))
	require.NoError(t, err)

	// Create Room Domain Object
	roomID := room_id.NewRoomId()
	roomName := "Test Room"
	password := "password123"
	qrCode := "qr-code-data"
	now := time.Now().UTC()
	lifespan := 3600

	r := room.HydrateRoom(
		roomID,
		roomName,
		password,
		qrCode,
		nil, // BoostCooldown
		now,
		uint32(lifespan),
		nil, // Songs
		[]user_id.UserId{userID},
	)

	// Mock Expectations
	hashedPassword := []byte("hashed_password")
	mockEncrypter.On("HashAndSalt", password).Return(hashedPassword, nil)

	// Act
	err = repo.CreateRoom(ctx, r, txx)
	require.NoError(t, err)

	// Assert Room Created
	var storedRoom struct {
		ID           uuid.UUID `db:"id"`
		Name         string    `db:"name"`
		Password     []byte    `db:"password"`
		QrCodeHash   []byte    `db:"qr_code_hash"`
		CreatedAtUtc time.Time `db:"created_at_utc"`
	}
	err = txx.GetContext(ctx, &storedRoom, "SELECT id, name, password, qr_code_hash, created_at_utc FROM rooms WHERE id = $1", uuid.UUID(roomID))
	require.NoError(t, err)
	assert.Equal(t, roomName, storedRoom.Name)
	assert.Equal(t, hashedPassword, storedRoom.Password)
	assert.Equal(t, []byte(qrCode), storedRoom.QrCodeHash)

	// Assert User Updated
	var userRoomID uuid.UUID
	err = txx.QueryRowContext(ctx, "SELECT room_id FROM users_room_data WHERE user_id = $1", uuid.UUID(userID)).Scan(&userRoomID)
	require.NoError(t, err)
	assert.Equal(t, uuid.UUID(roomID), userRoomID)

	// Assert User Role
	var role string
	err = txx.QueryRowContext(ctx, "SELECT role FROM users_room_data WHERE user_id = $1 AND room_id = $2", uuid.UUID(userID), uuid.UUID(roomID)).Scan(&role)
	require.NoError(t, err)
	assert.Equal(t, "host", role)
}

func TestRoomRepositoryGetRoomByUserId(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := room_repository.NewRoomRepository(mockEncrypter)

	userToTest := seeder.SeedData.Users[0]
	userID := userToTest.Id()

	// Act
	roomDao, err := repo.GetRoomByUserId(ctx, user_id.UserId(userID), txx)
	require.NoError(t, err)

	// Assert
	usersRoomId := userToTest.RoomId()
	usersRoom, _ := slice_extensions.
		GetRoomById(
			seeder.SeedData.Rooms,
			*usersRoomId,
		)

	assert.Equal(t, usersRoom.Name(), roomDao.Name)
	usersRole := userToTest.Role().String()
	assert.Equal(t, *usersRole, roomDao.UserRole)
	songsInRoom := usersRoom.EnqueuedSongs()
	assert.Len(t, songsInRoom, len(roomDao.SongDaos))

	songDao := roomDao.SongDaos[0]
	songFromDb, _ := slice_extensions.GetEnqueuedSongById(seeder.SeedData.Songs, enqueued_song_id.EnqueuedSongId(songDao.Id))
	assert.Equal(t, songFromDb.SongData().Title(), songDao.Title)
	assert.Equal(t, songFromDb.SongData().Author(), songDao.Author)
	userThatAddedSong, _ := slice_extensions.GetUserById(seeder.SeedData.Users, songFromDb.AddedBy())
	assert.Equal(t, userThatAddedSong.FullName().String(), songDao.AddedBy) // Concat name + surname
	assert.Equal(t, vote_status.Upvoted.String(), songDao.VoteStatus)
}

func TestRoomRepositoryCheckUserMembership(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)

	repo := room_repository.NewRoomRepository(mockEncrypter)

	// Setup User with Room
	userID1 := uuid.New()
	roomID := uuid.New()
	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'User', 'One')`, userID1)
	require.NoError(t, err)
	_, err = txx.ExecContext(ctx, `INSERT INTO users_room_data (room_id, user_id, role) VALUES ($1, $2, 'member')`, roomID, userID1)
	require.NoError(t, err)

	// Setup User without Room
	userID2 := uuid.New()
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'User', 'Two')`, userID2)
	require.NoError(t, err)

	// Act & Assert
	exists1 := repo.CheckUserMembership(ctx, user_id.UserId(userID1), txx)
	assert.True(t, exists1)

	exists2 := repo.CheckUserMembership(ctx, user_id.UserId(userID2), txx)
	assert.False(t, exists2)
}

func TestRoomRepositoryGetRoomIdByNameAndPassword(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := room_repository.NewRoomRepository(mockEncrypter)
	userIdToLeaveRoomFrom := seeder.SeedData.Rooms[0].Members()[0]
	err := repo.LeaveRoom(ctx, userIdToLeaveRoomFrom, txx)

	require.NoError(t, err)
	var isUserInRoom bool
	err = txx.GetContext(ctx, &isUserInRoom, "SELECT EXISTS(SELECT 1 FROM users_room_data WHERE user_id = $1)", uuid.UUID(userIdToLeaveRoomFrom))
	require.NoError(t, err)
	assert.False(t, isUserInRoom)
}

func TestRoomRepositoryLeaveRoom(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)
	repo := room_repository.NewRoomRepository(mockEncrypter)

	// Setup Data
	userID := uuid.New()
	roomID := uuid.New()

	_, err := txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Alice', 'Wonder')`, userID)
	require.NoError(t, err)

	// Act
	err = repo.LeaveRoom(ctx, user_id.UserId(userID), txx)
	require.NoError(t, err)

	// Assert
	var userDb user_dao.UserDao
	err = txx.GetContext(ctx, &userDb, "SELECT * FROM users WHERE id = $1", userID)
	require.NoError(t, err)
	assert.Nil(t, userDb.RoomId)
}

func TestRoomRepositoryJoinRoomById(t *testing.T) {
	txx,
		ctx,
		mockEncrypter := setupMocks(t)

	repo := room_repository.NewRoomRepository(mockEncrypter)

	// Setup: create a user and a room
	userID := uuid.New()
	roomID := uuid.New()

	_, err := txx.ExecContext(ctx, `INSERT INTO users (id, name, surname) VALUES ($1, 'Join', 'User')`, userID)
	require.NoError(t, err)

	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Join Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	// Act
	err = repo.JoinRoomById(ctx, user_id.UserId(userID), room_id.RoomId(roomID), txx)
	require.NoError(t, err)

	// Assert: user is member of the room with default role 'member'
	var role string
	err = txx.QueryRowContext(ctx, `SELECT role FROM users_room_data WHERE user_id = $1 AND room_id = $2`, userID, roomID).Scan(&role)
	require.NoError(t, err)
	assert.Equal(t, "member", role)
}
