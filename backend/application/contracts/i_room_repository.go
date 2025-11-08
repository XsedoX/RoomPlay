package contracts

import (
	"context"

	"xsedox.com/main/domain/room"
)

type IRoomRepository interface {
	Create(ctx context.Context, roomParam *room.Room, queryer IQueryer) error
}
