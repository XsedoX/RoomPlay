package dependecies

import "xsedox.com/main/infrastructure/persistance"

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/room"
	"xsedox.com/main/presentation/handlers"
)

type RoomDependencies struct {
	roomRepository    *persistance.RoomRepository
	createRoomHandler *room.CreateCommandHandler
	roomHandler       *handlers.RoomHandler
	unitOfWork        *persistance.UnitOfWork
}

func NewRoomDependencies(db *sqlx.DB) *RoomDependencies {
	roomRepository := persistance.NewRoomRepository(db)
	unitOfWork := persistance.NewUnitOfWork(db)
	createRoomHandler := room.NewCreateCommandHandler(roomRepository, unitOfWork)
	roomHandler := handlers.NewRoomHandler(createRoomHandler)

	return &RoomDependencies{
		roomRepository:    roomRepository,
		createRoomHandler: createRoomHandler,
		roomHandler:       roomHandler,
		unitOfWork:        unitOfWork,
	}
}
func (roomDeps *RoomDependencies) GetRoomHandler() *handlers.RoomHandler {
	return roomDeps.roomHandler
}
