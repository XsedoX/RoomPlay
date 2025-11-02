package shared

import "github.com/google/uuid"

type RoomId uuid.UUID

func (id *RoomId) String() *string {
	if id == nil {
		return nil
	}
	result := uuid.UUID(*id).String()
	return &result
}
func (id *RoomId) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	rId := uuid.UUID(*id)
	return &rId
}
