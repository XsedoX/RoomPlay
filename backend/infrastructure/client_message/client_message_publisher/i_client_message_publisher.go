package client_message_publisher

import (
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/i_client_message_handler"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/connection_id"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
)

type IClientMessagePublisher interface {
	RegisterHandler(messageName websocket_action.WebSocketIncomingAction, handler i_client_message_handler.IClientMessageHandler)
	Publish(rawEnvelope []byte, userId user_id.UserId, connectionId connection_id.ConnectionId)
	Run()
}
