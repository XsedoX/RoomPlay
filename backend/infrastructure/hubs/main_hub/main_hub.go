package main_hub

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/room_hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/client_room_request"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/room_broadcast_request"
	"github.com/gorilla/websocket"
)

type Hub struct {
	roomHubs map[room_id.RoomId]*room_hub.RoomHub

	// Register requests from the clients.
	register chan *client_room_request.ClientRoomRequest

	// Unregister requests from clients.
	unregister chan *client_room_request.ClientRoomRequest

	roomBroadcast chan *room_broadcast_request.RoomBroadcastRequest

	appContext context.Context
}

func NewHub(appContext context.Context) *Hub {
	return &Hub{
		register:      make(chan *client_room_request.ClientRoomRequest, 100),
		unregister:    make(chan *client_room_request.ClientRoomRequest, 100),
		roomHubs:      make(map[room_id.RoomId]*room_hub.RoomHub),
		roomBroadcast: make(chan *room_broadcast_request.RoomBroadcastRequest, 100),
		appContext:    appContext,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case joinRequest := <-h.register:
			roomHub, ok := h.roomHubs[joinRequest.RoomId]
			if !ok {
				roomHub = room_hub.NewRoomHub(h.appContext)
				go roomHub.Run()
			}
			h.roomHubs[joinRequest.RoomId] = roomHub
			roomHub.RegisterClient(joinRequest.Client)
		case leaveRequest := <-h.unregister:
			roomHub, ok := h.roomHubs[leaveRequest.RoomId]
			if ok {
				roomHub.UnregisterClient(leaveRequest.Client)
			}
		case broadcastRequest := <-h.roomBroadcast:
			roomHub, ok := h.roomHubs[broadcastRequest.RoomId]
			if ok {
				roomHub.Broadcast(broadcastRequest.Payload)
			}
		case <-h.appContext.Done():
			return
		}
	}
}

func (h *Hub) BroadcastToRoom(bradcastRequest *room_broadcast_request.RoomBroadcastRequest) {
	h.roomBroadcast <- bradcastRequest
}

func (h *Hub) RegisterClientToRoom(joinRequest *client_room_request.ClientRoomRequest) {
	h.register <- joinRequest
}

func NewWebSocketUpgrader() *websocket.Upgrader {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &upgrader
}
