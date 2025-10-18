package room

import (
	"xsedox.com/domain/entities"
)

type CreateCommandHandler struct {
	roomRepository IRepository
}

func NewCreateCommandHandler(roomRepository IRepository) *CreateCommandHandler {
	return &CreateCommandHandler{
		roomRepository: roomRepository,
	}
}
func (h *CreateCommandHandler) Handle(cmd CreateCommand) error {
	room, err := entities.NewRoom(cmd.RoomName, cmd.RoomPassword, entities.UserId(cmd.UserId))
	if err != nil {
		return err
	}
	return h.roomRepository.Create(room)
}
