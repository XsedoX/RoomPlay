package join

import "github.com/google/uuid"

type RoomCommand struct {
	RoomName     string    `json:"roomName" validate:"required"`
	RoomPassword string    `json:"roomPassword" validate:"required"`
	UserId       uuid.UUID `json:"userId" validate:"required"`
}
