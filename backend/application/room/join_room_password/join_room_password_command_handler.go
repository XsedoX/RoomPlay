package join_room_password

import (
	"context"

	"xsedox.com/main/application"
	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/application/room/contracts"
)

type JoinRoomPasswordCommandHandler struct {
	roomRepository contracts.IRoomRepository
	unitOfWork     contracts2.IUnitOfWork
}

func NewJoinRoomPasswordCommandHandler(roomRepository contracts.IRoomRepository, unitOfWork contracts2.IUnitOfWork) *JoinRoomPasswordCommandHandler {
	return &JoinRoomPasswordCommandHandler{roomRepository: roomRepository, unitOfWork: unitOfWork}
}

func (handler *JoinRoomPasswordCommandHandler) Handle(ctx context.Context, command *JoinRoomPasswordCommand) error {
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return application.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		roomId, getRoomIdErr := handler.roomRepository.GetRoomIdByNameAndPassword(ctx, command.RoomName, command.RoomPassword, handler.unitOfWork.GetQueryer())
		if getRoomIdErr != nil {
			return custom_errors.NewCustomError("JoinRoomPasswordCommandHandler.GetRoomIdByNameAndPassword",
				"Room with given name and password does not exist.",
				getRoomIdErr,
				custom_errors.NotFound)
		}
		joinRoomErr := handler.roomRepository.JoinRoomById(ctx, *userId, *roomId, handler.unitOfWork.GetQueryer())
		if joinRoomErr != nil {
			return custom_errors.NewCustomError("JoinRoomPasswordCommandHandler.JoinRoomById",
				"Something went wrong while joining the room.",
				joinRoomErr,
				custom_errors.Unexpected)
		}
		return nil
	})
	return err
}
