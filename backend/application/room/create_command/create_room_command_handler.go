package create_command

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
)

type CreateRoomCommandHandler struct {
	roomRepository contracts.IRoomRepository
	unitOfWork     contracts.IUnitOfWork
	encrypter      contracts.IEncrypter
}

func NewCreateRoomCommandHandler(roomRepository contracts.IRoomRepository,
	unitOfWork contracts.IUnitOfWork,
	encrypter contracts.IEncrypter) *CreateRoomCommandHandler {
	return &CreateRoomCommandHandler{
		roomRepository: roomRepository,
		unitOfWork:     unitOfWork,
		encrypter:      encrypter,
	}
}

func (handler CreateRoomCommandHandler) Handle(ctx context.Context, command *CreateRoomCommand) (*shared.RoomId, error) {
	var response shared.RoomId
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		qrCode := handler.encrypter.NewEncryptionKey()
		userId, ok := application.GetUserIdFromContext(ctx)
		if !ok {
			return applicationErrors.NewApplicationError("CreateRoomCommandHandler.GetUserIdFromContext",
				application.MissingUserIdInContextErrorMessage,
				nil,
				applicationErrors.Forbidden)
		}
		roomInstance := room.NewRoom(command.RoomName, command.RoomPassword, string(qrCode), *userId)
		err := handler.roomRepository.Create(ctx, roomInstance, handler.unitOfWork.GetQueryer())
		if err != nil {
			return applicationErrors.NewApplicationError("CreateRoomCommandHandler.Create",
				"Problem with creating a room.",
				err,
				applicationErrors.Unexpected)
		}
		response = roomInstance.Id()
		return nil
	})
	return &response, err
}
