package i_client_message_handler

import (
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/client_message_publisher_request"
)

type IClientMessageHandler interface {
	HandleMessage(client_message_publisher_request.ClientMessagePublisherRequest)
}
