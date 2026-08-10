package song_enqueued_websocket_event

import (
	"context"
	"log"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room/events"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

type SongEnqueuedWebsocketEventHandler struct {
	hub            hub.IHub
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
	appContext     context.Context
}

func NewSongEnqueuedWebsocketEventHandler(
	hub hub.IHub,
	roomRepo i_room_repository.IRoomRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	appContext context.Context,
) *SongEnqueuedWebsocketEventHandler {
	return &SongEnqueuedWebsocketEventHandler{
		hub:            hub,
		roomRepository: roomRepo,
		unitOfWork:     unitOfWork,
		appContext:     appContext,
	}
}

func (h *SongEnqueuedWebsocketEventHandler) Handle(event shared.IDomainEvent) {
	concreteEvent, ok := event.(*events.SongEnqueuedEvent)
	if !ok {
		log.Printf("Received event of unexpected type: %T", event)
		return
	}

	id := concreteEvent.EnqueuedSongId()

	addedBy, err := h.roomRepository.GetEnqueuedSongAddedByValueByRoomIdEnqueuedSongId(
		h.appContext,
		concreteEvent.RoomId(),
		id,
		h.unitOfWork.GetQueryer(h.appContext),
	)
	if err != nil {
		// log an error
		log.Printf("Error retrieving addedBy for enqueued song: %v", err)
		return
	}

	responseDto := SongEnqueuedWebsocketEventResponse{
		Id:            *id.ToUuid(),
		Author:        concreteEvent.Author(),
		Title:         concreteEvent.Title(),
		Votes:         concreteEvent.Votes(),
		AlbumCoverUrl: concreteEvent.AlbumCoverUrl(),
		State:         concreteEvent.EnqueuedSongState().String(),
		VoteStatus:    concreteEvent.Status().String(),
		AddedBy:       addedBy,
	}
	jsonPatchBuilder := response.NewJsonPatchResponseBuilder()
	idString := id.ToUuid().String()
	jsonPatch := jsonPatchBuilder.
		Add(idString, responseDto).
		Build(websocket_action.EnqueuedSongsPatchActionName)

	h.hub.BroadcastToRoom(&hub.RoomBroadcastRequest{
		RoomId:  concreteEvent.RoomId(),
		Payload: jsonPatch,
	})
}
