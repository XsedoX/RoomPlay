package song_enqueued_websocket_event

import (
	"context"
	"encoding/json"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room/events"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/i_hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/room_broadcast_request"
)

const songEnqueuedWebsocketAction = "song_enqueued"

type SongEnqueuedWebsocketEventHandler struct {
	hub            i_hub.IHub
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
	appContext     context.Context
}

func NewSongEnqueuedWebsocketEventHandler(
	hub i_hub.IHub,
	roomRepo i_room_repository.IRoomRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	appContext context.Context,
) *SongEnqueuedWebsocketEventHandler {
	return &SongEnqueuedWebsocketEventHandler{
		hub:            hub,
		roomRepository: roomRepo,
		unitOfWork:     unitOfWork,
	}
}

func (h *SongEnqueuedWebsocketEventHandler) Handle(event shared.IDomainEvent) {
	concreteEvent, _ := event.(*events.SongEnqueuedEvent)

	id := concreteEvent.EnqueuedSongId()

	addedBy, _ := h.roomRepository.GetEnqueuedSongAddedByValueByRoomIdEnqueuedSongId(
		h.appContext,
		concreteEvent.RoomId(),
		id,
		h.unitOfWork.GetQueryer(),
	)

	dto := SongEnqueuedWebsocketEventResponse{
		Id:            *id.ToUuid(),
		Author:        concreteEvent.Author(),
		Title:         concreteEvent.Title(),
		Votes:         concreteEvent.Votes(),
		AlbumCoverUrl: concreteEvent.AlbumCoverUrl(),
		State:         concreteEvent.EnqueuedSongState().String(),
		VoteStatus:    concreteEvent.Status().String(),
		AddedBy:       addedBy,
		Action:        songEnqueuedWebsocketAction,
	}
	payload, _ := json.Marshal(dto)

	h.hub.BroadcastToRoom(&room_broadcast_request.RoomBroadcastRequest{
		RoomId:  concreteEvent.RoomId(),
		Payload: payload,
	})
}
