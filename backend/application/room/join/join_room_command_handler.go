package join

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
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

func (h *RoomCommandHandler) Handle(ctx context.Context, cmd *RoomCommand) *applicationErrors.ApplicationError {
	err := h.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		room1, err := room.CreateRoom(cmd.RoomName, cmd.RoomPassword, shared.UserId(cmd.UserId))
		if err != nil {
			return err
		}
		return h.roomRepository.Join(ctx, room1, h.unitOfWork.GetQueryer())
	})
	respErr := applicationErrors.NewApplicationError("problem with transaction", err, applicationErrors.ErrInfrastructure)
	return respErr
}
