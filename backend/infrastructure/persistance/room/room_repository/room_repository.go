package room_repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_song_dao"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/google/uuid"
)

type RoomRepository struct {
	encrypter i_encrypter.IEncrypter
}

func NewRoomRepository(encrypter i_encrypter.IEncrypter) *RoomRepository {
	return &RoomRepository{
		encrypter: encrypter,
	}
}

func (repository *RoomRepository) GetEnqueuedSongAddedByValueByRoomIdEnqueuedSongId(ctx context.Context, roomId room_id.RoomId, enqueuedSongId enqueued_song_id.EnqueuedSongId, queryer i_queryer.IQueryer) (string, error) {
	var addedBy string
	err := queryer.GetContext(ctx,
		&addedBy,
		`
select concat(u."name" , ' ', u.surname) as added_by
from enqueued_songs es
join rooms r on es.room_id = r.id
join users u on es.added_by = u.id
where es.id = $1 and r.id = $2;
		`,
		enqueuedSongId.ToUuid(),
		roomId.ToUuid(),
	)
	if err != nil {
		return "", err
	}
	return addedBy, nil
}

func (repository *RoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer i_queryer.IQueryer) error {
	roomId := roomParam.Id()
	userId := roomParam.Members()[0]
	encryptedQrCode, err := repository.encrypter.Encrypt(roomParam.QrCode())
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
		roomParam.Password(),
		encryptedQrCode,
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

func (repository *RoomRepository) UpdateRoom(ctx context.Context, roomParam *room.Room, queryer i_queryer.IQueryer) error {
	roomId := roomParam.Id()

	// --- rooms table (aggregate root own columns) ---
	// Nullable column: boost_cooldown_seconds -> *uint16 (NULL when nil)
	encryptedQrCode, qrErr := repository.encrypter.Encrypt(roomParam.QrCode())
	if qrErr != nil {
		return qrErr
	}
	_, roomUpdateErr := queryer.ExecContext(ctx,
		`
		UPDATE rooms
		SET name = $1,
		    password = $2::bytea,
		    qr_code_hash = $3::bytea,
		    boost_cooldown_seconds = $4::smallint,  -- nullable -> *uint16 / NULL
		    created_at_utc = $5,
		    lifespan_seconds = $6
		WHERE id = $7::uuid;
		`,
		roomParam.Name(),
		roomParam.Password(),
		encryptedQrCode,
		roomParam.BoostCooldownSeconds(), // nil -> NULL
		roomParam.CreatedAtUtc(),
		roomParam.LifespanSeconds(),
		roomId.ToUuid(),
	)
	if roomUpdateErr != nil {
		return roomUpdateErr
	}

	// --- members (users_room_data) ---
	// Note: room.Members() does not carry roles; the host/role handling is
	// preserved below by keeping existing host rows, then syncing members.
	_, delMembersErr := queryer.ExecContext(ctx,
		`
		DELETE FROM users_room_data
		WHERE room_id = $1::uuid
		  AND role <> 'host';  -- keep the host row intact
		`,
		roomId.ToUuid(),
	)
	if delMembersErr != nil {
		return delMembersErr
	}
	for _, member := range roomParam.Members() {
		_, insMemberErr := queryer.ExecContext(ctx,
			`
			INSERT INTO users_room_data (room_id, user_id, role)
			VALUES ($1::uuid, $2::uuid, 'member')
			ON CONFLICT (user_id, room_id) DO NOTHING;
			`,
			roomId.ToUuid(),
			member.ToUuid(),
		)
		if insMemberErr != nil {
			return insMemberErr
		}
	}

	// --- banned_users ---
	_, delBannedErr := queryer.ExecContext(ctx,
		`DELETE FROM banned_users WHERE room_id = $1::uuid;`,
		roomId.ToUuid(),
	)
	if delBannedErr != nil {
		return delBannedErr
	}
	for _, bannedUser := range roomParam.BannedUsers() {
		_, insBannedErr := queryer.ExecContext(ctx,
			`
			INSERT INTO banned_users (room_id, user_id)
			VALUES ($1::uuid, $2::uuid);
			`,
			roomId.ToUuid(),
			bannedUser.ToUuid(),
		)
		if insBannedErr != nil {
			return insBannedErr
		}
	}

	upsertSongQuery := `
	insert into songs (title, author, isrc, id)
	values ($1, $2, $3, gen_random_uuid())
	on conflict (title, author, isrc) do update
	set title = excluded.title,
	author = excluded.author,
	isrc = excluded.isrc
	returning id::uuid;
	`
	upsertSongExternalDataQuery := `
			INSERT INTO songs_external_data (song_id, length_seconds, album_cover_url, url, music_provider)
			VALUES ($1::uuid, $2::smallint, $3, $4, $5)
			ON CONFLICT (song_id) DO UPDATE
			  SET length_seconds = EXCLUDED.length_seconds,
			      album_cover_url = EXCLUDED.album_cover_url,
			      url = EXCLUDED.url,
			      music_provider = EXCLUDED.music_provider;
	`
	// --- scheduledSong ---
	if roomParam.ScheduledSong() != nil {
		scheduledSong := roomParam.ScheduledSong()
		sd := scheduledSong.SongData()
		var scheduledSongSongId uuid.UUID
		insSongErr := queryer.GetContext(ctx,
			&scheduledSongSongId,
			upsertSongQuery,
			scheduledSong.SongData().Title(),
			scheduledSong.SongData().Author(),
			scheduledSong.SongData().Isrc(), // *string -> songs.isrc is NULLABLE (varchar(12) | nullable)
		)
		if insSongErr != nil {
			return insSongErr
		}
		_, insExtErr := queryer.ExecContext(ctx,
			upsertSongExternalDataQuery,
			scheduledSongSongId,
			sd.LengthSeconds(),
			sd.AlbumCoverUrl(),
			sd.Url(),
			sd.MusicProvider().String(),
		)
		if insExtErr != nil {
			return insExtErr
		}
		_, insSchedErr := queryer.ExecContext(ctx,
			`
		INSERT INTO scheduled_songs (room_id, song_id, scheduled_at_utc)
		VALUES ($1::uuid, $2::uuid, $3)
			on conflict (room_id) do update
			set song_id = excluded.song_id,
			    scheduled_at_utc = excluded.scheduled_at_utc;
		`,
			roomId.ToUuid(),
			scheduledSongSongId,
			scheduledSong.ScheduledAtUtc(),
		)
		if insSchedErr != nil {
			return insSchedErr
		}
	}

	// --- enqueued_songs ---
	_, delEnqueuedErr := queryer.ExecContext(ctx,
		`
		DELETE FROM enqueued_songs WHERE room_id = $1::uuid;`,
		roomId.ToUuid(),
	)
	if delEnqueuedErr != nil {
		return delEnqueuedErr
	}
	enqueuedSongs := roomParam.AllSongs()
	for _, enqueuedSong := range enqueuedSongs {
		songData := enqueuedSong.SongData()
		var enqueuedSongSongId uuid.UUID
		insSongErr := queryer.GetContext(ctx,
			&enqueuedSongSongId,
			upsertSongQuery,
			songData.Title(),
			songData.Author(),
			songData.Isrc(), // *string -> songs.isrc is NULLABLE (varchar(12) | nullable)
		)
		if insSongErr != nil {
			return insSongErr
		}
		_, insExtErr := queryer.ExecContext(ctx,
			upsertSongExternalDataQuery,
			enqueuedSongSongId,
			songData.LengthSeconds(),
			songData.AlbumCoverUrl(),
			songData.Url(),
			songData.MusicProvider().String(),
		)
		if insExtErr != nil {
			return insExtErr
		}
		enqueuedSongId := enqueuedSong.Id()
		_, insEnqueuedErr := queryer.ExecContext(ctx,
			`
			insert into enqueued_songs (id, room_id, song_id, added_at_utc, started_at_utc, state, added_by)
			values ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6, $7::uuid);
			`,
			enqueuedSongId.ToUuid(),
			roomId.ToUuid(),
			enqueuedSongSongId,
			enqueuedSong.AddedAtUtc(),
			enqueuedSong.StartedAtUtc(),
			enqueuedSong.State().String(),
			enqueuedSong.AddedBy().ToUuid(),
		)
		if insEnqueuedErr != nil {
			return insEnqueuedErr
		}
	}

	// --- default_playlist (default_playlists) ---
	// room_id is PRIMARY KEY -> upsert on conflict (room_id).
	defaultPlaylist := roomParam.DefaultPlaylist()
	_, delDpErr := queryer.ExecContext(ctx,
		`DELETE FROM default_playlists WHERE room_id = $1::uuid;`,
		roomId.ToUuid(),
	)
	if delDpErr != nil {
		return delDpErr
	}
	if defaultPlaylist != nil {
		_, insDpErr := queryer.ExecContext(ctx,
			`
			INSERT INTO default_playlists (external_id, user_id, songs_amount, playlist_title, room_id)
			VALUES ($1, $2::uuid, $3::smallint, $4, $5::uuid);
			`,
			defaultPlaylist.ExternalId(),
			defaultPlaylist.UserId().ToUuid(),
			uint16(defaultPlaylist.SongsAmount()), // uint -> smallint (NOT NULL)
			defaultPlaylist.Title(),
			roomId.ToUuid(),
		)
		if insDpErr != nil {
			return insDpErr
		}
	}

	return nil
}

func (repository *RoomRepository) GetRoomById(ctx context.Context, roomId room_id.RoomId, queryer i_queryer.IQueryer) (*room.Room, error) {
	type roomDao struct {
		Id                   uuid.UUID    `db:"id"`
		Name                 string       `db:"name"`
		Password             []byte       `db:"password"`
		QrCodeHash           []byte       `db:"qr_code_hash"`
		BoostCooldownSeconds *uint16      `db:"boost_cooldown_seconds"`
		CreatedAtUtc         sql.NullTime `db:"created_at_utc"`
		LifespanSeconds      uint32       `db:"lifespan_seconds"`
	}
	var roomDaoInstance roomDao
	getRoomErr := queryer.GetContext(ctx,
		&roomDaoInstance,
		`
SELECT id,
       name,
       password,                -- bytea (stored hash)
       qr_code_hash,            -- bytea (stored hash)
       boost_cooldown_seconds,  -- smallint | nullable -> *uint16
       created_at_utc,
       lifespan_seconds
FROM rooms
WHERE id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getRoomErr != nil {
		return nil, getRoomErr
	}

	var members []uuid.UUID
	getMembersErr := queryer.SelectContext(ctx,
		&members,
		`
SELECT user_id
FROM users_room_data
WHERE room_id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getMembersErr != nil {
		return nil, getMembersErr
	}
	membersAsUserIds := make([]user_id.UserId, 0, len(members))
	for _, member := range members {
		membersAsUserIds = append(membersAsUserIds, user_id.UserId(member))
	}

	var bannedUsers []uuid.UUID
	getBannedUsersErr := queryer.SelectContext(ctx,
		&bannedUsers,
		`
SELECT user_id
FROM banned_users
WHERE room_id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getBannedUsersErr != nil {
		return nil, getBannedUsersErr
	}
	bannedUsersAsUserIds := make([]user_id.UserId, 0, len(bannedUsers))
	for _, bannedUser := range bannedUsers {
		bannedUsersAsUserIds = append(bannedUsersAsUserIds, user_id.UserId(bannedUser))
	}

	type enqueuedSongDao struct {
		Id            uuid.UUID  `db:"id"`
		AddedAtUtc    time.Time  `db:"added_at_utc"`
		StartedAtUtc  *time.Time `db:"started_at_utc"`
		State         string     `db:"state"`
		AddedBy       uuid.UUID  `db:"added_by"`
		Votes         int8       `db:"votes"`
		Title         string     `db:"title"`
		Author        string     `db:"author"`
		Isrc          *string    `db:"isrc"`
		Url           *string    `db:"url"`
		AlbumCoverUrl *string    `db:"album_cover_url"`
		LengthSeconds *uint16    `db:"length_seconds"`
		MusicProvider *string    `db:"music_provider"`
	}
	var enqueuedSongsInstances []enqueuedSongDao
	getEnqueuedSongsErr := queryer.SelectContext(ctx,
		&enqueuedSongsInstances,
		`
SELECT es.id,
       es.added_at_utc,
       es.started_at_utc,
       es.state,
       es.added_by,
       COALESCE(v.votes, 0)::int8          AS votes,
       s.title, s.author, s.isrc,
       sed.url, sed.album_cover_url,
       sed.length_seconds, sed.music_provider
FROM enqueued_songs es
JOIN songs s              ON s.id = es.song_id
LEFT JOIN songs_external_data sed
       ON sed.song_id = es.song_id
LEFT JOIN (
    SELECT enqueued_song_id,
           SUM(CASE WHEN vote_status = 'upvoted'   THEN 1
                    WHEN vote_status = 'downvoted' THEN -1
                    ELSE 0 END)::int8 AS votes
    FROM users_votes
    GROUP BY enqueued_song_id
) v ON v.enqueued_song_id = es.id
WHERE es.room_id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getEnqueuedSongsErr != nil {
		return nil, getEnqueuedSongsErr
	}

	type scheduledSongDao struct {
		ScheduledAtUtc time.Time `db:"scheduled_at_utc"`
		Title          string    `db:"title"`
		Author         string    `db:"author"`
		Isrc           *string   `db:"isrc"`
		Url            *string   `db:"url"`
		AlbumCoverUrl  *string   `db:"album_cover_url"`
		LengthSeconds  *uint16   `db:"length_seconds"`
		MusicProvider  *string   `db:"music_provider"`
	}
	var scheduledSongsInstances []scheduledSongDao
	getScheduledSongsErr := queryer.SelectContext(ctx,
		&scheduledSongsInstances,
		`
SELECT ss.scheduled_at_utc,
       s.title, s.author, s.isrc,
       sed.url, sed.album_cover_url,
       sed.length_seconds, sed.music_provider
FROM scheduled_songs ss
JOIN songs s              ON s.id = ss.song_id
LEFT JOIN songs_external_data sed
       ON sed.song_id = ss.song_id
WHERE ss.room_id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getScheduledSongsErr != nil {
		return nil, getScheduledSongsErr
	}

	type defaultPlaylistDao struct {
		ExternalId    string `db:"external_id"`
		UserId        string `db:"user_id"`
		SongsAmount   int    `db:"songs_amount"`
		PlaylistTitle string `db:"playlist_title"`
	}
	var defaultPlaylistsInstances []defaultPlaylistDao
	getDefaultPlaylistsErr := queryer.SelectContext(ctx,
		&defaultPlaylistsInstances,
		`
SELECT external_id, user_id, songs_amount, playlist_title
FROM default_playlists
WHERE room_id = $1::uuid;
		`,
		roomId.ToUuid(),
	)
	if getDefaultPlaylistsErr != nil {
		return nil, getDefaultPlaylistsErr
	}

	enqueuedSongs := make([]enqueued_song.EnqueuedSong, 0, len(enqueuedSongsInstances))
	for _, enqueuedSongDaoInstance := range enqueuedSongsInstances {
		songData := song_data.HydrateSongData(
			*enqueuedSongDaoInstance.Url,
			enqueuedSongDaoInstance.Title,
			enqueuedSongDaoInstance.Author,
			*enqueuedSongDaoInstance.AlbumCoverUrl,
			*enqueuedSongDaoInstance.LengthSeconds,
			*music_provider.ParseMusicProvider(*enqueuedSongDaoInstance.MusicProvider),
			enqueuedSongDaoInstance.Isrc,
		)
		addedBy := user_id.UserId(enqueuedSongDaoInstance.AddedBy)
		enqueuedSongId := enqueued_song_id.EnqueuedSongId(enqueuedSongDaoInstance.Id)
		enqueuedSongInstance := enqueued_song.HydrateEnqueuedSong(
			enqueuedSongId,
			*songData,
			enqueuedSongDaoInstance.AddedAtUtc,
			enqueuedSongDaoInstance.StartedAtUtc,
			*enqueued_song_state.ParseSongState(enqueuedSongDaoInstance.State),
			enqueuedSongDaoInstance.Votes,
			addedBy,
		)
		enqueuedSongs = append(enqueuedSongs, *enqueuedSongInstance)
	}

	qrCodeString, qrDecryptErr := repository.encrypter.Decrypt(roomDaoInstance.QrCodeHash)
	if qrDecryptErr != nil {
		return nil, qrDecryptErr
	}
	roomInstance := room.HydrateRoom(
		room_id.RoomId(roomDaoInstance.Id),
		roomDaoInstance.Name,
		roomDaoInstance.Password,
		qrCodeString,
		roomDaoInstance.BoostCooldownSeconds,
		roomDaoInstance.CreatedAtUtc.Time,
		roomDaoInstance.LifespanSeconds,
		enqueuedSongs,
		membersAsUserIds,
	)
	return roomInstance, nil
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
