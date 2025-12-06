package persistance

import (
	"context"
	"database/sql"

	"xsedox.com/main/application/contracts"
	daos2 "xsedox.com/main/application/room/get_room/daos"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type RoomRepository struct {
	encrypter contracts.IEncrypter
}

func NewRoomRepository(encrypter contracts.IEncrypter) *RoomRepository {
	return &RoomRepository{
		encrypter: encrypter,
	}
}
func (repository *RoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer contracts.IQueryer) error {
	roomId := roomParam.Id()
	userId := roomParam.Members()[0]
	hashedPassword, err := repository.encrypter.HashAndSalt(roomParam.Password())
	if err != nil {
		return err
	}
	_, addRoomErr := queryer.ExecContext(ctx,
		`INSERT INTO rooms VALUES ($1::uuid, $2, $3::bytea, $4::bytea, $5, $6, $7);`,
		roomId.ToUuid(),
		roomParam.Name(),
		hashedPassword,
		[]byte(roomParam.QrCode()),
		roomParam.BoostCooldownSeconds(),
		roomParam.CreatedAtUtc().UTC(),
		roomParam.LifespanSeconds(),
	)
	if addRoomErr != nil {
		return addRoomErr
	}

	role := user.Host
	_, usersUpdateErr := queryer.ExecContext(ctx,
		`INSERT INTO users_room_data (room_id, user_id, role) VALUES ($1, $2, $3)`,
		roomId.ToUuid(),
		userId.ToUuid(),
		role.String(),
	)
	if usersUpdateErr != nil {
		return usersUpdateErr
	}
	return nil
}
func (repository *RoomRepository) GetRoomByUserId(ctx context.Context, userId user.Id, queryer contracts.IQueryer) (*daos2.GetRoomDao, error) {
	var getRoomDaoInstance daos2.GetRoomDao
	getRoomErr := queryer.GetContext(ctx,
		&getRoomDaoInstance,
		`
SELECT rooms.name,
	   rooms.qr_code_hash,
	   users_room_data.boost_used_at_utc,
	   rooms.boost_cooldown_seconds,
	   songs.title AS playing_song_title,
	   songs.author AS playing_song_author,
	   enqueued_songs.started_at_utc AS playing_song_started_at_utc,
	   songs.length_seconds AS playing_song_length_seconds,
	   users_room_data.role
FROM rooms
JOIN users_room_data ON users_room_data.room_id = rooms.id
LEFT JOIN enqueued_songs ON enqueued_songs.room_id = rooms.id AND enqueued_songs.state = 'playing'
LEFT JOIN songs ON songs.id = enqueued_songs.song_id
WHERE users_room_data.user_id = $1::uuid 
LIMIT 1;
`, userId.ToUuid())
	if getRoomErr != nil {
		return nil, getRoomErr
	}

	getRoomsSongDaoInstances := make([]daos2.GetRoomSongDao, 0)
	getRoomSongsErr := queryer.SelectContext(ctx,
		&getRoomsSongDaoInstances,
		`
SELECT enqueued_songs.id,
       songs.title,
       songs.author,
       CONCAT(users_for_added_by.name, ' ', users_for_added_by.surname) AS added_by,
       enqueued_songs.state,
       COALESCE(users_votes.vote_status, 'not_voted') AS vote_status,
       enqueued_songs.votes,
       songs.album_cover_url
FROM enqueued_songs
         JOIN songs ON enqueued_songs.song_id = songs.id
         JOIN rooms ON enqueued_songs.room_id = rooms.id
         JOIN users_room_data ON users_room_data.room_id = rooms.id
         JOIN users AS users_for_added_by ON users_for_added_by.id = enqueued_songs.added_by
         LEFT JOIN users_votes ON users_room_data.user_id = users_votes.user_id AND enqueued_songs.id = users_votes.enqueued_song_id
WHERE users_room_data.user_id = $1;
`, userId.ToUuid())
	if getRoomSongsErr != nil {
		return nil, getRoomSongsErr
	}
	getRoomDaoInstance.SongDaos = getRoomsSongDaoInstances
	return &getRoomDaoInstance, nil
}
func (repository *RoomRepository) CheckUserMembership(ctx context.Context, userId user.Id, queryer contracts.IQueryer) bool {
	var response bool
	err := queryer.GetContext(ctx, &response, `
		SELECT CASE 
		    WHEN EXISTS (
		        SELECT 1
		        FROM users_room_data
		        WHERE user_id=$1
			)
		    THEN true 
		    ELSE false
		END`, userId.ToUuid())
	if err != nil {
		return false
	}
	return response
}
func (repository *RoomRepository) LeaveRoom(ctx context.Context, id user.Id, queryer contracts.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		`DELETE FROM users_room_data WHERE user_id=$1`,
		id.ToUuid())
	return err
}
func (repository *RoomRepository) JoinRoomById(ctx context.Context, userId user.Id, roomId shared.RoomId, queryer contracts.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		`INSERT INTO users_room_data (user_id, room_id) VALUES ($1, $2)`, userId.ToUuid(), roomId.ToUuid())
	return err
}
func (repository *RoomRepository) GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer contracts.IQueryer) (*shared.RoomId, error) {
	type roomWithNamePasswordAndId struct {
		roomId       shared.RoomId `db:"id"`
		roomPassword []byte        `db:"password"`
	}
	var rooms []roomWithNamePasswordAndId
	getRoomsErr := queryer.SelectContext(ctx,
		&rooms,
		`SELECT id, password FROM rooms WHERE name=$1;`, roomName)
	if getRoomsErr != nil {
		return nil, getRoomsErr
	}
	if len(rooms) == 0 {
		return nil, sql.ErrNoRows
	}
	for _, roomDb := range rooms {
		if repository.encrypter.Verify(roomPassword, roomDb.roomPassword) {
			return &roomDb.roomId, nil
		}
	}
	return nil, sql.ErrNoRows
}
