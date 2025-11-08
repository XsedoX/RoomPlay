package room

import (
	"github.com/google/uuid"
)

type SongId uuid.UUID

func (id *SongId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}
func (id *SongId) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	concreteId := uuid.UUID(*id)
	return &concreteId
}
func ParseSongId(s string) *SongId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := SongId(uuidResult)
	return &result
}
