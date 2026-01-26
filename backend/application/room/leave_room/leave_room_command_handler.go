package leave_room

import (
	"context"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts"
)

type LeaveRoomCommandHandler struct {
	roomRepository room_contracts.IRoomRepository
	unitOfWork     application_contracts.IUnitOfWork
}

func (l LeaveRoomCommandHandler) Handle(ctx context.Context, _ *LeaveRoomCommand) error {
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return application.NewMissingUserIdInContextError
	}
	err := l.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		err := l.roomRepository.LeaveRoom(ctx, *userId, l.unitOfWork.GetQueryer())
		if err != nil {
			return customerrors.NewCustomError("LeaveRoomCommandHandler.LeaveRoom",
				"Couldn't leave room.",
				err,
				customerrors.Unexpected)
		}
		return nil
	})
	return err
}

func NewLeaveRoomCommandHandler(roomRepository room_contracts.IRoomRepository, unitOfWork application_contracts.IUnitOfWork) *LeaveRoomCommandHandler {
	return &LeaveRoomCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}
