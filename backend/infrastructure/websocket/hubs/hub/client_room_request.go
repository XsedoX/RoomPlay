package hub

import (
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
)

type ClientRoomRequest struct {
	RoomId room_id.RoomId
	Client *Client
}
