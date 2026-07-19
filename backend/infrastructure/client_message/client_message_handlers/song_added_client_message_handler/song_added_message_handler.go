package song_added_client_message_handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	enquque_song_command "github.com/XsedoX/RoomPlay/application/song/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"
)

const SongAddedClientMessageName = "song_added"

type SongAddedClientMessageHandler struct {
	commandHandler i_command_handler.ICommandHandler[*enquque_song_command.EnqueueSongCommand]
}

func NewSongAddedClientMessageHandler(
	commandHandler i_command_handler.ICommandHandler[*enquque_song_command.EnqueueSongCommand],
) *SongAddedClientMessageHandler {
	return &SongAddedClientMessageHandler{
		commandHandler: commandHandler,
	}
}

func (handler *SongAddedClientMessageHandler) HandleMessage(envelope client_message_envelope.ClientMessageEnvelope) {
	var command enquque_song_command.EnqueueSongCommand
	err := json.Unmarshal(envelope.Payload, &command)
	if err != nil {
		log.Printf("Failed to unmarshal song added message: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = handler.commandHandler.Handle(ctx, &command)
	if err != nil {
		log.Printf("Failed to handle song added command: %s", err)
		return
	}
}
