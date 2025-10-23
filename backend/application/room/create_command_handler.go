package room

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/domain/entities"
)

type CreateCommandHandler struct {
	roomRepository IRepository
	unitOfWork     application.IUnitOfWork
}

func NewCreateCommandHandler(roomRepository IRepository, uow application.IUnitOfWork) *CreateCommandHandler {
	return &CreateCommandHandler{
		roomRepository: roomRepository,
		unitOfWork:     uow,
	}
}

func (h *CreateCommandHandler) Handle(ctx context.Context, cmd *CreateCommand) error {
	return h.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		room, err := entities.CreateRoom(cmd.RoomName, cmd.RoomPassword, entities.UserId(cmd.UserId))
		if err != nil {
			return err
		}
		return h.roomRepository.Create(ctx, room)
	})
}
