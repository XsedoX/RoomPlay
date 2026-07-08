package join_room_password_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
)

type JoinRoomPasswordCommandHandler struct {
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
}

func NewJoinRoomPasswordCommandHandler(roomRepository i_room_repository.IRoomRepository, unitOfWork i_unit_of_work.IUnitOfWork) *JoinRoomPasswordCommandHandler {
	return &JoinRoomPasswordCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}

func (handler *JoinRoomPasswordCommandHandler) Handle(ctx context.Context, command *join_room_password_command.JoinRoomPasswordCommand) error {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return application_helpers.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		roomId, getRoomIdErr := handler.roomRepository.GetRoomIdByNameAndPassword(ctx, command.RoomName, command.RoomPassword, handler.unitOfWork.GetQueryer())
		if getRoomIdErr != nil {
			return application_error.NewApplicationError("JoinRoomPasswordCommandHandler.GetRoomIdByNameAndPassword",
				"Room with given name and password does not exist.",
				getRoomIdErr,
				application_error_type.NotFound)
		}
		joinRoomErr := handler.roomRepository.JoinRoomById(ctx, *userId, *roomId, handler.unitOfWork.GetQueryer())
		if joinRoomErr != nil {
			return application_error.NewApplicationError("JoinRoomPasswordCommandHandler.JoinRoomById",
				"Something went wrong while joining the room.",
				joinRoomErr,
				application_error_type.Unexpected)
		}
		return nil
	})
	return err
}
