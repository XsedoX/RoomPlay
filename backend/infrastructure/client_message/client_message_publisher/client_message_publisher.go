package client_message_publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/i_client_message_handler"
)

type (
	ClientMessageName      string
	ClientMessagePublisher struct {
		clientMessageHandlers map[ClientMessageName]i_client_message_handler.IClientMessageHandler

		messages chan client_message_envelope.ClientMessageEnvelope

		applicationContext context.Context
	}
)

func NewClientMessagePublisher(applicationContext context.Context) *ClientMessagePublisher {
	return &ClientMessagePublisher{
		clientMessageHandlers: make(map[ClientMessageName]i_client_message_handler.IClientMessageHandler),
		messages:              make(chan client_message_envelope.ClientMessageEnvelope, 100),
		applicationContext:    applicationContext,
	}
}

func (p *ClientMessagePublisher) RegisterHandler(messageName ClientMessageName, handler i_client_message_handler.IClientMessageHandler) {
	p.clientMessageHandlers[messageName] = handler
}

func (p *ClientMessagePublisher) Publish(rawEnvelope []byte, userId user_id.UserId) {
	var envelope client_message_envelope.ClientMessageEnvelope
	err := json.Unmarshal(rawEnvelope, &envelope)
	if err != nil {
		log.Printf("Failed to unmarshal client message envelope: %s", err)
		return
	}
	envelope.UserId = userId
	p.messages <- envelope
}

func (p *ClientMessagePublisher) Run() {
	for {
		select {
		case message := <-p.messages:
			handler, ok := p.clientMessageHandlers[ClientMessageName(message.ActionName)]
			if !ok {
				continue
			}
			handler.HandleMessage(message)
		case <-p.applicationContext.Done():
			return
		}
	}
}
