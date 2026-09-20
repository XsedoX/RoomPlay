package hub

import "github.com/XsedoX/RoomPlay/domain/room/room_id"

// payload should be websocket response either failed or success
type RoomBroadcastRequest struct {
	RoomId  room_id.RoomId
	Payload any
}
