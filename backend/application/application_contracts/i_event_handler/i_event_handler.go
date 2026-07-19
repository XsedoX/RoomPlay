package i_event_handler

import "github.com/XsedoX/RoomPlay/domain/shared"

type IEventHandler interface {
	Handle(event shared.IDomainEvent)
}
