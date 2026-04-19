package enqueued_song_id

import "github.com/google/uuid"

type EnqueuedSongId uuid.UUID

func (id *EnqueuedSongId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}

func (id *EnqueuedSongId) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	concreteId := uuid.UUID(*id)
	return &concreteId
}

func ParseEnqueuedSongId(s string) *EnqueuedSongId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := EnqueuedSongId(uuidResult)
	return &result
}

func NewEnqueuedSongId() EnqueuedSongId {
	result := EnqueuedSongId(uuid.New())
	return result
}

func IdsEqual(id1, id2 *EnqueuedSongId) bool {
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
