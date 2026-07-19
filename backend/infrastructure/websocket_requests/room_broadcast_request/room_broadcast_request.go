package room_broadcast_request

import "github.com/XsedoX/RoomPlay/domain/room/room_id"

type RoomBroadcastRequest struct {
	RoomId  room_id.RoomId
	Payload []byte
}
