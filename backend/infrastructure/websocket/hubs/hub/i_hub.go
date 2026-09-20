package hub

import (
	"github.com/gorilla/websocket"
)

type IHub interface {
	BroadcastToRoom(roomBroadcastRequest *RoomBroadcastRequest)
	RegisterClientToRoom(clientRoomRequest *ClientRoomRequest)
	NewWebSocketUpgrader() *websocket.Upgrader
	UnregisterClient(clientRoomRequest *ClientRoomRequest)
	BroadcastToClientByConnectionId(broadcastRequest *ConnectionIdBroadcastRequest)
}
