package users_votes_dao

import (
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/google/uuid"
)

type UsersVotesDao struct {
	UserId         uuid.UUID                             `json:"user_id"`
	EnqueuedSongId uuid.UUID                             `json:"enqueued_song_id"`
	State          enqueued_song_state.EnqueuedSongState `json:"state"`
}
