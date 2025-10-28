package shared

import "github.com/google/uuid"

type RoomId uuid.UUID

func (id RoomId) String() string {
	return uuid.UUID(id).String()
}
