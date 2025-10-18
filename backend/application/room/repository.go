package room

import (
	"xsedox.com/domain/entities"
)

type IRepository interface {
	Create(room *entities.Room) error
}
