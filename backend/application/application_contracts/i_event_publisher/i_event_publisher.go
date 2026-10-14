package i_event_publisher

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_event_handler"
	"github.com/XsedoX/RoomPlay/domain/shared"
)

type IEventPublisher interface {
	Publish(event shared.IDomainEvent)
	Register(eventName shared.EventName, handler i_event_handler.IEventHandler)
	Start()
}
