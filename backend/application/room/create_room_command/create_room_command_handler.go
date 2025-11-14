package create_room_command

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts2 "xsedox.com/main/application/room/contracts"
	"xsedox.com/main/domain/room"
)

type CreateRoomCommandHandler struct {
	roomRepository contracts2.IRoomRepository
	unitOfWork     contracts.IUnitOfWork
	encrypter      contracts.IEncrypter
}

func NewCreateRoomCommandHandler(roomRepository contracts2.IRoomRepository,
	unitOfWork contracts.IUnitOfWork,
	encrypter contracts.IEncrypter) *CreateRoomCommandHandler {
	return &CreateRoomCommandHandler{
		roomRepository: roomRepository,
		unitOfWork:     unitOfWork,
		encrypter:      encrypter,
	}
}

func (handler CreateRoomCommandHandler) Handle(ctx context.Context, command *CreateRoomCommand) error {
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return application.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		qrCode := handler.encrypter.NewEncryptionKey()
		roomInstance := room.NewRoom(command.RoomName, command.RoomPassword, string(qrCode), *userId)
		err := handler.roomRepository.CreateRoom(ctx, roomInstance, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("CreateRoomCommandHandler.CreateRoom",
				"Problem with creating a room.",
				err,
				custom_errors.Unexpected)
		}
		return nil
	})
	return err
}
