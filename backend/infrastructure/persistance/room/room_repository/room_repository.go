package room_repository

import (
	"context"
	"database/sql"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_song_dao"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
)

type RoomRepository struct {
	encrypter i_encrypter.IEncrypter
}

func NewRoomRepository(encrypter i_encrypter.IEncrypter) *RoomRepository {
	return &RoomRepository{
		encrypter: encrypter,
	}
}

func (repository *RoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer i_queryer.IQueryer) error {
	roomId := roomParam.Id()
	userId := roomParam.Members()[0]
	hashedPassword, err := repository.encrypter.HashAndSalt(roomParam.Password())
	if err != nil {
		return err
	}
	_, addRoomErr := queryer.ExecContext(ctx,
		`
		INSERT INTO rooms
		VALUES 
		(
			$1::uuid, $2, $3::bytea, $4::bytea, $5, $6, $7
		);
		`,
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

	role := user_role.Host
	_, usersUpdateErr := queryer.ExecContext(ctx,
		`
		INSERT INTO users_room_data
		(
			room_id, user_id, role
		) 
		VALUES 
		(
			$1, $2, $3
		)
		`,
		roomId.ToUuid(),
		userId.ToUuid(),
		role.String(),
	)
	if usersUpdateErr != nil {
		return usersUpdateErr
	}
	return nil
}

func (repository *RoomRepository) GetRoomByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*get_room_dao.GetRoomDao, error) {
	var getRoomDaoInstance get_room_dao.GetRoomDao
	getRoomErr := queryer.GetContext(ctx,
		&getRoomDaoInstance,
		`
		select rooms.name,
	   rooms.qr_code_hash,
	   users_room_data.boost_used_at_utc,
	   rooms.boost_cooldown_seconds,
	   songs.title AS playing_song_title,
	   songs.author AS playing_song_author,
	   enqueued_songs.started_at_utc AS playing_song_started_at_utc,
	   songs_external_data.length_seconds AS playing_song_length_seconds,
	   users_room_data.role
FROM rooms
JOIN users_room_data ON users_room_data.room_id = rooms.id
LEFT JOIN enqueued_songs ON enqueued_songs.room_id = rooms.id AND enqueued_songs.state = 'playing'
left join users_room_data admin_data on admin_data.room_id = rooms.id and admin_data.role = 'host'
left join users_external_credentials admin_external_data on admin_data.user_id = admin_external_data.user_id 
LEFT JOIN songs_external_data ON songs_external_data.song_id = enqueued_songs.song_id and songs_external_data."music_provider" = admin_external_data."music_provider" 
left join songs on enqueued_songs.song_id = songs.id
WHERE users_room_data.user_id = $1::uuid 
LIMIT 1;
		`, userId.ToUuid())
	if getRoomErr != nil {
		return nil, getRoomErr
	}

	getRoomsSongDaoInstances := make([]get_room_song_dao.GetRoomSongDao, 0)
	getRoomSongsErr := queryer.SelectContext(ctx,
		&getRoomsSongDaoInstances,
		`
SELECT enqueued_songs.id,
       songs.title,
       songs.author,
       CONCAT(users_for_added_by.name, ' ', users_for_added_by.surname) AS added_by,
       enqueued_songs.state,
       COALESCE(users_votes.vote_status, 'not_voted') AS vote_status,
	   COALESCE(enqueued_songs_votes.value, 0) AS votes,
       songs_external_data.album_cover_url
FROM enqueued_songs
			JOIN songs ON enqueued_songs.song_id = songs.id
			JOIN rooms ON enqueued_songs.room_id = rooms.id
			JOIN users_room_data ON users_room_data.room_id = rooms.id
			JOIN users AS users_for_added_by ON users_for_added_by.id = enqueued_songs.added_by
		  left join songs_external_data on songs_external_data.song_id = enqueued_songs.song_id
			LEFT JOIN users_votes ON users_room_data.user_id = users_votes.user_id AND enqueued_songs.id = users_votes.enqueued_song_id
			LEFT JOIN (
			    SELECT enqueued_song_id,
			           SUM(CASE 
			                 WHEN vote_status = 'upvoted' THEN 1 
			                 WHEN vote_status = 'downvoted' THEN -1 
			                 ELSE 0 
			               END) AS value
			    FROM users_votes
			    GROUP BY enqueued_song_id
			) AS enqueued_songs_votes ON enqueued_songs.id = enqueued_songs_votes.enqueued_song_id
WHERE users_room_data.user_id = $1 and enqueued_songs.state != 'playing';
`, userId.ToUuid())
	if getRoomSongsErr != nil {
		return nil, getRoomSongsErr
	}
	getRoomDaoInstance.SongDaos = getRoomsSongDaoInstances
	return &getRoomDaoInstance, nil
}

func (repository *RoomRepository) CheckUserMembership(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) bool {
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

func (repository *RoomRepository) LeaveRoom(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		`DELETE FROM users_room_data WHERE user_id=$1`,
		id.ToUuid())
	return err
}

func (repository *RoomRepository) JoinRoomById(ctx context.Context, userId user_id.UserId, roomId room_id.RoomId, queryer i_queryer.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		`INSERT INTO users_room_data (user_id, room_id) VALUES ($1, $2)`, userId.ToUuid(), roomId.ToUuid())
	return err
}

func (repository *RoomRepository) GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer i_queryer.IQueryer) (*room_id.RoomId, error) {
	type roomWithNamePasswordAndId struct {
		RoomId       string `db:"id"`
		RoomPassword []byte `db:"password"`
	}
	var rooms []roomWithNamePasswordAndId
	getRoomsErr := queryer.SelectContext(ctx,
		&rooms,
		`SELECT id, password::bytea FROM rooms WHERE name=$1;`, roomName)
	if getRoomsErr != nil {
		return nil, getRoomsErr
	}
	if len(rooms) == 0 {
		return nil, sql.ErrNoRows
	}
	for _, roomDb := range rooms {
		if repository.encrypter.Verify(roomPassword, roomDb.RoomPassword) {
			return room_id.ParseRoomId(roomDb.RoomId), nil
		}
	}
	return nil, sql.ErrNoRows
}
