package enqueue_song_command_handler

import (
	"context"

	enquque_song_command "github.com/XsedoX/RoomPlay/application/song/enqueue_song/enqueue_song_command"
)

type EnqueueSongCommandHandler struct{}

func NewEnqueueSongCommandHandler() *EnqueueSongCommandHandler {
	return &EnqueueSongCommandHandler{}
}

func (e *EnqueueSongCommandHandler) Handle(context context.Context, command *enquque_song_command.EnqueueSongCommand) error {
	// Implement the logic to handle the EnqueueSongCommand here
	return nil
}

