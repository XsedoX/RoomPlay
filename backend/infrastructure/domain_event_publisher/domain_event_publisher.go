package domain_event_publisher

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_event_handler"
	"github.com/XsedoX/RoomPlay/domain/shared"
)

type DommainEventPublisher struct {
	eventQueue         chan shared.IDomainEvent
	eventHandlers      map[shared.EventName][]i_event_handler.IEventHandler
	applicationContext context.Context
}

func NewDomainEventPublisher(appContext context.Context) *DommainEventPublisher {
	return &DommainEventPublisher{
		eventQueue:         make(chan shared.IDomainEvent, 100),
		eventHandlers:      make(map[shared.EventName][]i_event_handler.IEventHandler),
		applicationContext: appContext,
	}
}

func (p *DommainEventPublisher) Publish(event shared.IDomainEvent) {
	p.eventQueue <- event
}

func (p *DommainEventPublisher) Register(eventName shared.EventName, handler i_event_handler.IEventHandler) {
	p.eventHandlers[eventName] = append(p.eventHandlers[eventName], handler)
}

func (p *DommainEventPublisher) Start() {
	for {
		select {
		case <-p.applicationContext.Done():
			return // goroutine exits
		case event := <-p.eventQueue:
			for _, handler := range p.eventHandlers[event.EventName()] {
				handler.Handle(event)
			}
		}
	}
}
