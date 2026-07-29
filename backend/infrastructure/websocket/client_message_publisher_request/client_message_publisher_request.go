package client_message_publisher_request

import (
	"encoding/json"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/connection_id"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
)

type ClientMessagePublisherRequest struct {
	ConnectionId connection_id.ConnectionId               `json:"connectionId"`
	UserId       user_id.UserId                           `json:"userId"`
	ActionName   websocket_action.WebSocketIncomingAction `json:"actionName"`
	Payload      json.RawMessage                          `json:"payload"`
}
