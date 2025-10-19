package room

import "github.com/google/uuid"

type CreateCommand struct {
	RoomName     string    `json:"roomName" validate:"required"`
	RoomPassword string    `json:"roomPassword" validate:"required"`
	UserId       uuid.UUID `json:"userId" validate:"required"`
}
