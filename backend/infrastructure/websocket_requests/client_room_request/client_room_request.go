package client_room_request

import (
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/room_hub"
)

type ClientRoomRequest struct {
	RoomId room_id.RoomId
	Client *room_hub.Client
}
