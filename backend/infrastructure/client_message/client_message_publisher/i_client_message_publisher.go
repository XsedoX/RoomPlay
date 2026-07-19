package client_message_publisher

import (
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/i_client_message_handler"
)

type IClientMessagePublisher interface {
	RegisterHandler(messageName ClientMessageName, handler i_client_message_handler.IClientMessageHandler)
	Publish(rawEnvelope []byte, userId user_id.UserId)
	Run()
}
