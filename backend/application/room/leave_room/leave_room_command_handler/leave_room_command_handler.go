package leave_room_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
)

type LeaveRoomCommandHandler struct {
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
}

func (l LeaveRoomCommandHandler) Handle(ctx context.Context, _ *leave_room_command.LeaveRoomCommand) error {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return application_helpers.NewMissingUserIdInContextError
	}
	err := l.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		err := l.roomRepository.LeaveRoom(ctx, *userId, l.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("LeaveRoomCommandHandler.LeaveRoom",
				"Couldn't leave room.",
				err,
				custom_error_type.Unexpected)
		}
		return nil
	})
	return err
}

func NewLeaveRoomCommandHandler(roomRepository i_room_repository.IRoomRepository, unitOfWork i_unit_of_work.IUnitOfWork) *LeaveRoomCommandHandler {
	return &LeaveRoomCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}
