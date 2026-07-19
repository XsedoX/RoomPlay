package i_hub

import (
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/client_room_request"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/room_broadcast_request"
)

type IHub interface {
	BroadcastToRoom(roomBroadcastRequest *room_broadcast_request.RoomBroadcastRequest)
	RegisterClientToRoom(clientRoomRequest *client_room_request.ClientRoomRequest)
}
