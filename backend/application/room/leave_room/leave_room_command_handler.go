package leave_room

import (
	"context"

	"github.com/XsedoX/RoomPlay/application"
	contracts2 "github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	contracts3 "github.com/XsedoX/RoomPlay/application/room/contracts"
)

type LeaveRoomCommandHandler struct {
	roomRepository contracts3.IRoomRepository
	unitOfWork     contracts2.IUnitOfWork
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

func NewLeaveRoomCommandHandler(roomRepository contracts3.IRoomRepository, unitOfWork contracts2.IUnitOfWork) *LeaveRoomCommandHandler {
	return &LeaveRoomCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}
