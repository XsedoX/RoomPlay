package persistance

import (
	"context"

	"xsedox.com/main/application/contracts"
	daos2 "xsedox.com/main/application/room/get_room_query/daos"
	"xsedox.com/main/domain/room"
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
func (rr *RoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer contracts.IQueryer) error {
	roomId := roomParam.Id()
	userId := roomParam.Members()[0]
	hashedPassword, err := rr.encrypter.HashAndSalt(roomParam.Password())
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

	_, usersUpdateErr := queryer.ExecContext(ctx,
		`UPDATE users SET room_id = $1::uuid WHERE id = $2::uuid;`,
		roomId.ToUuid(),
		userId.ToUuid(),
	)
	if usersUpdateErr != nil {
		return usersUpdateErr
	}

	host := user.Host
	_, usersRolesErr := queryer.ExecContext(ctx,
		`INSERT INTO users_roles VALUES ($1, $2, $3);`,
		roomId.ToUuid(),
		userId.ToUuid(),
		host.String(),
	)
	if usersRolesErr != nil {
		return usersRolesErr
	}
	return nil
}
func (rr *RoomRepository) GetRoomByUserId(ctx context.Context, userId user.Id, queryer contracts.IQueryer) (*daos2.GetRoomDao, error) {
	var getRoomDaoInstance daos2.GetRoomDao
	getRoomErr := queryer.GetContext(ctx,
		&getRoomDaoInstance,
		`
SELECT rooms.name,
	   rooms.qr_code_hash,
	   boosts.used_at_utc,
	   songs.title AS playing_song_title,
	   songs.author AS playing_song_author,
	   enqueued_songs.started_at_utc AS playing_song_started_at_utc,
	   songs.length_seconds AS playing_song_length_seconds,
	   users_roles.role
FROM rooms
JOIN users ON users.room_id = rooms.id
LEFT JOIN boosts ON boosts.room_id = rooms.id
				AND boosts.user_id = users.id
JOIN users_roles ON users.id = users_roles.user_id
				AND users_roles.room_id = rooms.id
LEFT JOIN enqueued_songs ON enqueued_songs.room_id = rooms.id AND enqueued_songs.state = 'playing'
LEFT JOIN songs ON songs.id = enqueued_songs.song_id
WHERE users.id = $1::uuid 
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
       enqueued_songs.votes,
       songs.album_cover_url
FROM enqueued_songs
JOIN songs ON enqueued_songs.song_id = songs.id
JOIN rooms ON enqueued_songs.room_id = rooms.id
JOIN users ON users.room_id = rooms.id
JOIN users AS users_for_added_by ON users_for_added_by.id = enqueued_songs.added_by
WHERE users.id = $1;
`, userId.ToUuid())
	if getRoomSongsErr != nil {
		return nil, getRoomSongsErr
	}
	getRoomDaoInstance.SongDaos = getRoomsSongDaoInstances
	return &getRoomDaoInstance, nil
}
func (rr *RoomRepository) CheckUserMembership(ctx context.Context, userId user.Id, queryer contracts.IQueryer) bool {
	var response bool
	err := queryer.GetContext(ctx, &response, `
		SELECT CASE 
		    WHEN EXISTS (
		        SELECT 1
		        FROM users
		        WHERE id=$1 AND room_id IS NOT NULL
			)
		    THEN true 
		    ELSE false
		END`, userId.ToUuid())
	if err != nil {
		return false
	}
	return response
}
