package leave_room

import (
	"context"

	"xsedox.com/main/application"
	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts3 "xsedox.com/main/application/room/contracts"
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
			return custom_errors.NewCustomError("LeaveRoomCommandHandler.LeaveRoom",
				"Couldn't leave room.",
				err,
				custom_errors.Unexpected)
		}
		return nil
	})
	return err
}

func NewLeaveRoomCommandHandler(roomRepository contracts3.IRoomRepository, unitOfWork contracts2.IUnitOfWork) *LeaveRoomCommandHandler {
	return &LeaveRoomCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}
