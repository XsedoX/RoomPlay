package persistance

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/test_helpers/infrustructure_test/authentication_mocks"
)

func TestRoomRepository_CreateRoom(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := NewRoomRepository(mockEncrypter)

	// Setup User
	userID := user.Id(uuid.New())
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-host', 'Host', 'User')`, uuid.UUID(userID))
	require.NoError(t, err)

	// Create Room Domain Object
	roomID := shared.RoomId(uuid.New())
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
		[]user.Id{userID},
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

	mockEncrypter.AssertExpectations(t)
}

func TestRoomRepository_GetRoomByUserId(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := NewRoomRepository(mockEncrypter)

	// Setup Data
	roomID := uuid.New()
	userID := uuid.New()
	songID := uuid.New()
	enqueuedSongID := uuid.New()

	// Insert Room
	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds, boost_cooldown_seconds) VALUES ($1, 'My Room', 'pass', 'qr', $2, 3600, 60)`, roomID, time.Now().UTC())
	require.NoError(t, err)

	// Insert User
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-user', 'Test', 'User')`, userID)
	require.NoError(t, err)

	// Insert Role
	_, err = txx.ExecContext(ctx, `INSERT INTO users_room_data (room_id, user_id, role) VALUES ($1, $2, 'member')`, roomID, userID)
	require.NoError(t, err)

	// Insert Song
	_, err = txx.ExecContext(ctx, `INSERT INTO songs (id, external_id, title, author, length_seconds, album_cover_url) VALUES ($1, 'ext-song', 'Song Title', 'Song Author', 180, 'url')`, songID)
	require.NoError(t, err)

	// Insert Enqueued Song
	_, err = txx.ExecContext(ctx, `INSERT INTO enqueued_songs (id, room_id, song_id, added_by, added_at_utc, state, votes) VALUES ($1, $2, $3, $4, $5, 'enqueued', 0)`, enqueuedSongID, roomID, songID, userID, time.Now().UTC())
	require.NoError(t, err)

	// Insert Vote (to test join)
	_, err = txx.ExecContext(ctx, `INSERT INTO users_votes (user_id, enqueued_song_id, vote_status) VALUES ($1, $2, 'upvoted')`, userID, enqueuedSongID)
	require.NoError(t, err)

	// Act
	roomDao, err := repo.GetRoomByUserId(ctx, user.Id(userID), txx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "My Room", roomDao.Name)
	assert.Equal(t, "member", roomDao.UserRole)
	assert.Len(t, roomDao.SongDaos, 1)

	songDao := roomDao.SongDaos[0]
	assert.Equal(t, "Song Title", songDao.Title)
	assert.Equal(t, "Song Author", songDao.Author)
	assert.Equal(t, "Test User", songDao.AddedBy) // Concat name + surname
	assert.Equal(t, "upvoted", songDao.VoteStatus)
}

func TestRoomRepository_CheckUserMembership(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()

	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := NewRoomRepository(mockEncrypter)

	// Setup User with Room
	userID1 := uuid.New()
	roomID := uuid.New()
	_, err = txx.ExecContext(ctx, `INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds) VALUES ($1, 'Room', 'pass', 'qr', $2, 3600)`, roomID, time.Now().UTC())
	require.NoError(t, err)
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-1', 'User', 'One')`, userID1)
	require.NoError(t, err)
	_, err = txx.ExecContext(ctx, `INSERT INTO users_room_data (room_id, user_id, role) VALUES ($1, $2, 'member')`, roomID, userID1)
	require.NoError(t, err)

	// Setup User without Room
	userID2 := uuid.New()
	_, err = txx.ExecContext(ctx, `INSERT INTO users (id, external_id, name, surname) VALUES ($1, 'ext-2', 'User', 'Two')`, userID2)
	require.NoError(t, err)

	// Act & Assert
	exists1 := repo.CheckUserMembership(ctx, user.Id(userID1), txx)
	assert.True(t, exists1)

	exists2 := repo.CheckUserMembership(ctx, user.Id(userID2), txx)
	assert.False(t, exists2)
}

func TestRoomRepository_GetRoomIdByNameAndPassword(t *testing.T) {
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	defer txx.Rollback()
	mockEncrypter := new(authentication_mocks.MockEncrypter)
	repo := NewRoomRepository(mockEncrypter)
}
