package room_hub

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type RoomHub struct {
	// Registered clients.
	clients map[*Client]bool

	userSessions map[user_id.UserId]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	disconnectUserSessions chan user_id.UserId

	appContext context.Context
}

func NewRoomHub(appContext context.Context) *RoomHub {
	return &RoomHub{
		broadcast:              make(chan []byte, 100),
		register:               make(chan *Client, 100),
		unregister:             make(chan *Client, 100),
		clients:                make(map[*Client]bool),
		userSessions:           make(map[user_id.UserId]map[*Client]bool),
		disconnectUserSessions: make(chan user_id.UserId, 100),
		appContext:             appContext,
	}
}

func (hub *RoomHub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			userId := client.UserId()
			if _, ok := hub.userSessions[userId]; !ok {
				hub.userSessions[userId] = make(map[*Client]bool)
			}
			hub.userSessions[userId][client] = true
		case client := <-hub.unregister:
			delete(hub.clients, client)
			userId := client.UserId()
			delete(hub.userSessions[userId], client)
		case message := <-hub.broadcast:
			for client := range hub.clients {
				client.SendMessage(message)
			}
		case userId := <-hub.disconnectUserSessions:
			if sessions, ok := hub.userSessions[userId]; ok {
				for client := range sessions {
					client.Close()
				}
				delete(hub.userSessions, userId)
			}
		case <-hub.appContext.Done():
			for client := range hub.clients {
				client.Close()
			}
			return
		}
	}
}

func (hub *RoomHub) RegisterClient(client *Client) {
	hub.register <- client
}

func (hub *RoomHub) UnregisterClient(client *Client) {
	hub.unregister <- client
}

func (hub *RoomHub) Broadcast(message []byte) {
	hub.broadcast <- message
}

func (hub *RoomHub) DisconnectUserSessions(userId user_id.UserId) {
	hub.disconnectUserSessions <- userId
}
