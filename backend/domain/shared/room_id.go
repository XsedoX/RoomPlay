package shared

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
func (id *RoomId) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	concreteId := uuid.UUID(*id)
	return &concreteId
}
func ParseRoomId(s string) *RoomId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := RoomId(uuidResult)
	return &result
}
