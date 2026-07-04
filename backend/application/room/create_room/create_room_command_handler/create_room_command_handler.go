package create_room_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room"
)

type CreateRoomCommandHandler struct {
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
	encrypter      i_encrypter.IEncrypter
}

func NewCreateRoomCommandHandler(roomRepository i_room_repository.IRoomRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
	encrypter i_encrypter.IEncrypter,
) *CreateRoomCommandHandler {
	return &CreateRoomCommandHandler{
		roomRepository: roomRepository,
		unitOfWork:     unitOfWork,
		encrypter:      encrypter,
	}
}

func (handler CreateRoomCommandHandler) Handle(ctx context.Context, command *create_room_command.CreateRoomCommand) error {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return application_helpers.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		qrCode := handler.encrypter.NewEncryptionKey()
		roomInstance, domainErr := room.NewRoom(command.RoomName, command.RoomPassword, string(qrCode), *userId)
		if domainErr != nil {
			return domainErr
		}
		err := handler.roomRepository.CreateRoom(ctx, roomInstance, handler.unitOfWork.GetQueryer())
		if err != nil {
			return custom_error.NewCustomError("CreateRoomCommandHandler.CreateRoom",
				"Problem with creating a room.",
				err,
				custom_error_type.Unexpected)
		}
		return nil
	})
	return err
}
