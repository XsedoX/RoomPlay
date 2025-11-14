package contracts

import (
	"context"

	"xsedox.com/main/domain/shared"
)

type IEventPublisher interface {
	Publish(ctx context.Context, event shared.IDomainEvent)
}
