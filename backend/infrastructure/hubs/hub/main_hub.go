package hub

import (
	"context"
	"net/http"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/connection_id"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[room_id.RoomId]map[*Client]bool

	userSessions map[user_id.UserId]map[*Client]bool

	clientConnections map[connection_id.ConnectionId]*Client

	// Register requests from the clients.
	register chan *ClientRoomRequest

	// Unregister requests from clients.
	unregister chan *ClientRoomRequest

	roomBroadcast chan *RoomBroadcastRequest

	clientBroadcast chan *ConnectionIdBroadcastRequest

	appContext context.Context

	configuration config.IConfiguration

	disconnectUserSessions chan user_id.UserId
}

func NewHub(
	appContext context.Context,
	configuration config.IConfiguration,
) *Hub {
	return &Hub{
		register:               make(chan *ClientRoomRequest, 100),
		unregister:             make(chan *ClientRoomRequest, 100),
		roomBroadcast:          make(chan *RoomBroadcastRequest, 100),
		clientBroadcast:        make(chan *ConnectionIdBroadcastRequest, 100),
		disconnectUserSessions: make(chan user_id.UserId, 100),
		rooms:                  make(map[room_id.RoomId]map[*Client]bool),
		userSessions:           make(map[user_id.UserId]map[*Client]bool),
		clientConnections:      make(map[connection_id.ConnectionId]*Client),
		appContext:             appContext,
		configuration:          configuration,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case joinRequest := <-h.register:
			roomsClients, ok := h.rooms[joinRequest.RoomId]
			if !ok {
				roomsClients = make(map[*Client]bool)
			}
			roomsClients[joinRequest.Client] = true
			h.rooms[joinRequest.RoomId] = roomsClients

			userSessions, ok := h.userSessions[joinRequest.Client.UserId()]
			if !ok {
				userSessions = make(map[*Client]bool)
			}
			userSessions[joinRequest.Client] = true
			h.userSessions[joinRequest.Client.UserId()] = userSessions

			h.clientConnections[joinRequest.Client.ConnectionId()] = joinRequest.Client

			go joinRequest.Client.writePump()
			go joinRequest.Client.readPump()

		case leaveRequest := <-h.unregister:
			h.removeClient(leaveRequest.Client)

		case broadcastRequest := <-h.roomBroadcast:
			roomsClients, ok := h.rooms[broadcastRequest.RoomId]
			if ok {
				for client := range roomsClients {
					client.sendMessage(broadcastRequest.Payload)
				}
			}

		case userId := <-h.disconnectUserSessions:
			userSessions, ok := h.userSessions[userId]
			if ok {
				for client := range userSessions {
					h.removeClient(client)
				}
			}

		case connectionIdBroadcastRequest := <-h.clientBroadcast:
			client, ok := h.clientConnections[connectionIdBroadcastRequest.ConnectionId]
			if ok {
				client.sendMessage(connectionIdBroadcastRequest.Payload)
			}

		case <-h.appContext.Done():
			return
		}
	}
}

func (h *Hub) BroadcastToClientByConnectionId(broadcastRequest *ConnectionIdBroadcastRequest) {
	h.clientBroadcast <- broadcastRequest
}

func (h *Hub) DisconnectUserSessions(userId user_id.UserId) {
	h.disconnectUserSessions <- userId
}

func (h *Hub) UnregisterClient(clientRoomRequest *ClientRoomRequest) {
	h.unregister <- clientRoomRequest
}

func (h *Hub) BroadcastToRoom(bradcastRequest *RoomBroadcastRequest) {
	h.roomBroadcast <- bradcastRequest
}

func (h *Hub) RegisterClientToRoom(joinRequest *ClientRoomRequest) {
	h.register <- joinRequest
}

func (h *Hub) NewWebSocketUpgrader() *websocket.Upgrader {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == h.configuration.Authentication().ClientOrigin
	}}
	return &upgrader
}

func (h *Hub) removeClient(client *Client) {
	registeredClient, exists := h.clientConnections[client.ConnectionId()]
	if !exists || registeredClient != client {
		return
	}

	delete(h.clientConnections, client.ConnectionId())

	if roomClients, exists := h.rooms[client.RoomId()]; exists {
		delete(roomClients, client)

		if len(roomClients) == 0 {
			delete(h.rooms, client.RoomId())
		}
	}

	if userSessions, exists := h.userSessions[client.UserId()]; exists {
		delete(userSessions, client)

		if len(userSessions) == 0 {
			delete(h.userSessions, client.UserId())
		}
	}

	client.shutdown()
}
