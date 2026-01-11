package daos

import (
	"github.com/google/uuid"
	"github.com/XsedoX/RoomPlay/domain/room"
)

type UsersVotesDao struct {
	UserId         uuid.UUID      `json:"user_id"`
	EnqueuedSongId uuid.UUID      `json:"enqueued_song_id"`
	State          room.SongState `json:"state"`
}
