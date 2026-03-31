package get_room_song_dao

import (
	"time"

	"github.com/google/uuid"
)

type GetRoomSongDao struct {
	Id            uuid.UUID  `db:"id"`
	Url           string     `db:"url"`
	Title         string     `db:"title"`
	Author        string     `db:"author"`
	LengthSeconds uint16     `db:"length_seconds"`
	AlbumCoverUrl string     `db:"album_cover_url"`
	AddedAtUtc    time.Time  `db:"added_at_utc"`
	StartedAtUtc  *time.Time `db:"started_at_utc"`
	State         string     `db:"state"`
	Votes         uint8      `db:"votes"`
	AddedBy       string     `db:"added_by"`
}
