package contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/shared"
)

type IEventPublisher interface {
	Publish(ctx context.Context, event shared.IDomainEvent)
}
