package song_enqueued_client_message_handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	enqueue_song_command "github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/client_message_publisher_request"
)

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

func (handler *SongEnqueuedClientMessageHandler) HandleMessage(request client_message_publisher_request.ClientMessagePublisherRequest) {
	var command enqueue_song_command.EnqueueSongCommand
	err := json.Unmarshal(request.Payload, &command)
	if err != nil {
		log.Printf("Failed to unmarshal song added message: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	userId := request.UserId
	ctxWithClaims := context.WithValue(ctx, application_helpers.IdClaimContextKeyName, &userId)
	defer cancel()

	err = handler.commandHandler.Handle(ctxWithClaims, &command)
	if err != nil {
		log.Printf("Failed to handle song added command: %s", err)
		return
	}
}
