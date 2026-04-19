package room_id

import (
	"github.com/google/uuid"
)

type RoomId uuid.UUID

func (id *RoomId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}

func (id RoomId) ToUuid() uuid.UUID {
	concreteId := uuid.UUID(id)
	return concreteId
}

func ParseRoomId(s string) *RoomId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := RoomId(uuidResult)
	return &result
}

func NewRoomId() RoomId {
	return RoomId(uuid.New())
}

func IdsEqual(id1, id2 *RoomId) bool {
	if id1 == nil && id2 == nil {
		return true
	}
	if id1 == nil || id2 == nil {
		return false
	}
	if *id1 == *id2 {
		return true
	}
	return false
}
