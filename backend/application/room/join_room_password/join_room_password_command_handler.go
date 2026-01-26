package join_room_password

import (
	"context"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts"
)

type JoinRoomPasswordCommandHandler struct {
	roomRepository room_contracts.IRoomRepository
	unitOfWork     application_contracts.IUnitOfWork
}

func NewJoinRoomPasswordCommandHandler(roomRepository room_contracts.IRoomRepository, unitOfWork application_contracts.IUnitOfWork) *JoinRoomPasswordCommandHandler {
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
			return customerrors.NewCustomError("JoinRoomPasswordCommandHandler.GetRoomIdByNameAndPassword",
				"Room with given name and password does not exist.",
				getRoomIdErr,
				customerrors.NotFound)
		}
		joinRoomErr := handler.roomRepository.JoinRoomById(ctx, *userId, *roomId, handler.unitOfWork.GetQueryer())
		if joinRoomErr != nil {
			return customerrors.NewCustomError("JoinRoomPasswordCommandHandler.JoinRoomById",
				"Something went wrong while joining the room.",
				joinRoomErr,
				customerrors.Unexpected)
		}
		return nil
	})
	return err
}
