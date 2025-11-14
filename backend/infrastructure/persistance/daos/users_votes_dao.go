package daos

import (
	"github.com/google/uuid"
	"xsedox.com/main/domain/room"
)

type UsersVotesDao struct {
	UserId         uuid.UUID      `json:"user_id"`
	EnqueuedSongId uuid.UUID      `json:"enqueued_song_id"`
	State          room.SongState `json:"state"`
}
