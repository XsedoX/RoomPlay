package i_client_message_handler

import "github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"

type IClientMessageHandler interface {
	HandleMessage(client_message_envelope.ClientMessageEnvelope)
}
