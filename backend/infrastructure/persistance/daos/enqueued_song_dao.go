package daos

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/google/uuid"
)

type EnqueuedSongDao struct {
	Id           uuid.UUID                             `db:"id"`
	RoomId       uuid.UUID                             `db:"room_id"`
	SongId       uuid.UUID                             `db:"song_id"`
	AddedBy      uuid.UUID                             `db:"added_by"`
	AddedAtUtc   time.Time                             `db:"added_at_utc"`
	StartedAtUtc *time.Time                            `db:"started_at_utc"`
	State        enqueued_song_state.EnqueuedSongState `db:"state"`
	Votes        uint                                  `db:"votes"`
}
