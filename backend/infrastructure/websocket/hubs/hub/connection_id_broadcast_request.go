package hub

import "github.com/XsedoX/RoomPlay/infrastructure/websocket/hubs/connection_id"

type ConnectionIdBroadcastRequest struct {
	ConnectionId connection_id.ConnectionId `json:"connectionId"`
	Payload      any                        `json:"payload"`
}
