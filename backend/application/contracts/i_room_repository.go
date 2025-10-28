package contracts

import (
	"context"

	"xsedox.com/main/domain/room"
)

type IRoomRepository interface {
	Join(ctx context.Context, room *room.Room, queryer IQueryer) error
}
