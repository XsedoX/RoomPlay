package song_enqueued_client_message_handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	enqueue_song_command "github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"
)

const SongEnqueuedClientMessageActionName = "song_added"

type SongEnqueuedClientMessageHandler struct {
	commandHandler i_command_handler.ICommandHandler[*enqueue_song_command.EnqueueSongCommand]
}

func NewSongEnqueuedClientMessageHandler(
	commandHandler i_command_handler.ICommandHandler[*enqueue_song_command.EnqueueSongCommand],
) *SongEnqueuedClientMessageHandler {
	return &SongEnqueuedClientMessageHandler{
		commandHandler: commandHandler,
	}
}

func (handler *SongEnqueuedClientMessageHandler) HandleMessage(envelope client_message_envelope.ClientMessageEnvelope) {
	var command enqueue_song_command.EnqueueSongCommand
	err := json.Unmarshal(envelope.Payload, &command)
	if err != nil {
		log.Printf("Failed to unmarshal song added message: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	userId := envelope.UserId
	ctxWithClaims := context.WithValue(ctx, user.IdClaimContextKeyName, &userId)
	defer cancel()

	err = handler.commandHandler.Handle(ctxWithClaims, &command)
	if err != nil {
		log.Printf("Failed to handle song added command: %s", err)
		return
	}
}
