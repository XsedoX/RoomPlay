package connection_id

import "github.com/google/uuid"

type ConnectionId uuid.UUID

func (id *ConnectionId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}

func (id ConnectionId) ToUuid() uuid.UUID {
	concreteId := uuid.UUID(id)
	return concreteId
}

func ParseUserId(s string) *ConnectionId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := ConnectionId(uuidResult)
	return &result
}

func New() ConnectionId {
	return ConnectionId(uuid.New())
}

func IdsEqual(id1, id2 *ConnectionId) bool {
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
