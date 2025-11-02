package join

import (
	"context"

	"xsedox.com/main/application/contracts"
)

type RoomCommandHandler struct {
	roomRepository contracts.IRoomRepository
	unitOfWork     contracts.IUnitOfWork
}

func NewRoomCommandHandler(roomRepository contracts.IRoomRepository, uow contracts.IUnitOfWork) *RoomCommandHandler {
	return &RoomCommandHandler{
		roomRepository: roomRepository,
		unitOfWork:     uow,
	}
}

func (h *RoomCommandHandler) Handle(ctx context.Context, cmd *RoomCommand) error {
	//TODO TO BE IMPLEMENTED
	return nil
}
