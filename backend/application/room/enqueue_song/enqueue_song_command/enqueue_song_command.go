package enqueue_song_command

import "github.com/google/uuid"

type EnqueueSongCommand struct {
	SongExternalId string    `json:"SongExternalId"`
	AddedBy        uuid.UUID `json:"addedBy"`
}
