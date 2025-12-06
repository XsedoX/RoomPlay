package daos

import "time"

type GetRoomDao struct {
	Name                     string     `db:"name"`
	QrCode                   []byte     `db:"qr_code_hash"`
	PlayingSongTitle         *string    `db:"playing_song_title"`
	PlayingSongAuthor        *string    `db:"playing_song_author"`
	PlayingSongStartedAtUtc  *time.Time `db:"playing_song_started_at_utc"`
	PlayingSongLengthSeconds *uint16    `db:"playing_song_length_seconds"`
	UserRole                 string     `db:"role"`
	BoostUsedAtUtc           *time.Time `db:"boost_used_at_utc"`
	BoostCooldownSeconds     *uint16    `db:"boost_cooldown_seconds"`
	SongDaos                 []GetRoomSongDao
}
