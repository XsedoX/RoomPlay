package song_enqueued_client_message_handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	enqueue_song_command "github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/client_message_publisher_request"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/presentation/response"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
)

type SongEnqueuedClientMessageHandler struct {
	commandHandler i_command_handler.ICommandHandler[*enqueue_song_command.EnqueueSongCommand]
	mainHub        hub.IHub
}

func NewSongEnqueuedClientMessageHandler(
	commandHandler i_command_handler.ICommandHandler[*enqueue_song_command.EnqueueSongCommand],
	mainHub hub.IHub,
) *SongEnqueuedClientMessageHandler {
	return &SongEnqueuedClientMessageHandler{
		commandHandler: commandHandler,
		mainHub:        mainHub,
	}
}

func (handler *SongEnqueuedClientMessageHandler) HandleMessage(request client_message_publisher_request.ClientMessagePublisherRequest) {
	var command enqueue_song_command.EnqueueSongCommand
	err := json.Unmarshal(request.Payload, &command)
	if err != nil {
		errResponse := response.WebSocketFailure{
			ActionName: websocket_action.SongEnqueuedErrorActionName,
			ProblemDetails: response.ProblemDetails{
				Type:        "SongEnqueuedClientMessageHandler.Decoding",
				Title:       "Failed to decode the song enqueued command.",
				Status:      http.StatusBadRequest,
				Description: err.Error(),
				Instance:    constants.ApiBasePath + room_controller.RoomBasePath + room_controller.WebSocketUpgradePath,
			},
		}
		handler.mainHub.BroadcastToClientByConnectionId(&hub.ConnectionIdBroadcastRequest{
			ConnectionId: request.ConnectionId,
			Payload:      errResponse,
		})
		log.Printf("Failed to unmarshal song added message: %s", err)
		return
	}
	validationErr := setup_validation.ValidatorInstance.Struct(command)
	if validationErr != nil {
		errResponse := response.WebSocketFailure{
			ActionName: websocket_action.SongEnqueuedErrorActionName,
			ProblemDetails: response.ProblemDetails{
				Type:        "SongEnqueuedClientMessageHandler.ValidationError",
				Title:       "Validation error occurred while processing the song enqueued command.",
				Status:      http.StatusUnprocessableEntity,
				Description: validationErr.Error(),
				Instance:    constants.ApiBasePath + room_controller.RoomBasePath + room_controller.WebSocketUpgradePath,
			},
		}
		handler.mainHub.BroadcastToClientByConnectionId(&hub.ConnectionIdBroadcastRequest{
			ConnectionId: request.ConnectionId,
			Payload:      errResponse,
		})
		log.Printf("Failed to validate song added command: %s", validationErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	userId := request.UserId
	ctxWithClaims := context.WithValue(ctx, application_helpers.IdClaimContextKeyName, &userId)
	defer cancel()

	err = handler.commandHandler.Handle(ctxWithClaims, &command)
	if err != nil {
		errResponse := response.WebSocketFailure{
			ActionName: websocket_action.SongEnqueuedErrorActionName,
			ProblemDetails: response.ProblemDetails{
				Type:        "SongEnqueuedClientMessageHandler.EnqueueSongCommandHandler",
				Title:       "Error occurred while processing the song enqueued command.",
				Status:      http.StatusInternalServerError,
				Description: err.Error(),
				Instance:    constants.ApiBasePath + room_controller.RoomBasePath + room_controller.WebSocketUpgradePath,
			},
		}
		handler.mainHub.BroadcastToClientByConnectionId(&hub.ConnectionIdBroadcastRequest{
			ConnectionId: request.ConnectionId,
			Payload:      errResponse,
		})
		log.Printf("Failed to validate song added command: %s", validationErr)
		return
	}
}
