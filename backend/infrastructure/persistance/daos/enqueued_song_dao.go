package daos

import (
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/room"
)

type EnqueuedSongDao struct {
	Id           uuid.UUID      `db:"id"`
	RoomId       uuid.UUID      `db:"room_id"`
	SongId       uuid.UUID      `db:"song_id"`
	AddedBy      uuid.UUID      `db:"added_by"`
	AddedAtUtc   time.Time      `db:"added_at_utc"`
	StartedAtUtc *time.Time     `db:"started_at_utc"`
	State        room.SongState `db:"state"`
	Votes        uint           `db:"votes"`
}
