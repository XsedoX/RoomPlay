package room

import (
	"context"

	"xsedox.com/main/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, room *entities.Room) error
}
