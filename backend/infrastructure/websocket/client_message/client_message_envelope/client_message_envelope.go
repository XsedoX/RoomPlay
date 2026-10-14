package client_message_envelope

import (
	"encoding/json"

	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
)

type ClientMessageEnvelope struct {
	ActionName websocket_action.WebSocketIncomingAction `json:"actionName"`
	Payload    json.RawMessage                          `json:"payload"`
}
