package room

import "github.com/google/uuid"

type CreateCommand struct {
	RoomName     string    `json:"roomName"`
	RoomPassword string    `json:"roomPassword"`
	UserId       uuid.UUID `json:"userId"`
}
