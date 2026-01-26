package create_room

import (
	"context"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts"
	"github.com/XsedoX/RoomPlay/domain/room"
)

type CreateRoomCommandHandler struct {
	roomRepository room_contracts.IRoomRepository
	unitOfWork     application_contracts.IUnitOfWork
	encrypter      application_contracts.IEncrypter
}

func NewCreateRoomCommandHandler(roomRepository room_contracts.IRoomRepository,
	unitOfWork application_contracts.IUnitOfWork,
	encrypter application_contracts.IEncrypter,
) *CreateRoomCommandHandler {
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
		roomInstance, domainErr := room.NewRoom(command.RoomName, command.RoomPassword, string(qrCode), *userId)
		if domainErr != nil {
			return domainErr
		}
		err := handler.roomRepository.CreateRoom(ctx, roomInstance, handler.unitOfWork.GetQueryer())
		if err != nil {
			return customerrors.NewCustomError("CreateRoomCommandHandler.CreateRoom",
				"Problem with creating a room.",
				err,
				customerrors.Unexpected)
		}
		return nil
	})
	return err
}
