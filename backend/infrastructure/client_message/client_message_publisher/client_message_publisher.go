package client_message_publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/i_client_message_handler"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/connection_id"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/client_message_publisher_request"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
)

type (
	ClientMessagePublisher struct {
		clientMessageHandlers map[websocket_action.WebSocketIncomingAction]i_client_message_handler.IClientMessageHandler

		messages chan client_message_publisher_request.ClientMessagePublisherRequest

		applicationContext context.Context
	}
)

func NewClientMessagePublisher(applicationContext context.Context) *ClientMessagePublisher {
	return &ClientMessagePublisher{
		clientMessageHandlers: make(map[websocket_action.WebSocketIncomingAction]i_client_message_handler.IClientMessageHandler),
		messages:              make(chan client_message_publisher_request.ClientMessagePublisherRequest, 100),
		applicationContext:    applicationContext,
	}
}

func (p *ClientMessagePublisher) RegisterHandler(messageName websocket_action.WebSocketIncomingAction, handler i_client_message_handler.IClientMessageHandler) {
	p.clientMessageHandlers[messageName] = handler
}

func (p *ClientMessagePublisher) Publish(rawEnvelope []byte, userId user_id.UserId, connectionId connection_id.ConnectionId) {
	var envelope client_message_envelope.ClientMessageEnvelope
	err := json.Unmarshal(rawEnvelope, &envelope)
	if err != nil {
		log.Printf("Failed to unmarshal client message envelope: %s", err)
		return
	}
	request := client_message_publisher_request.ClientMessagePublisherRequest{
		ActionName:   envelope.ActionName,
		Payload:      envelope.Payload,
		UserId:       userId,
		ConnectionId: connectionId,
	}
	p.messages <- request
}

func (p *ClientMessagePublisher) Run() {
	for {
		select {
		case message := <-p.messages:
			handler, ok := p.clientMessageHandlers[websocket_action.WebSocketIncomingAction(message.ActionName)]
			if !ok {
				log.Printf("No handler registered for message: %s", message.ActionName)
				continue
			}
			handler.HandleMessage(message)
		case <-p.applicationContext.Done():
			return
		}
	}
}
