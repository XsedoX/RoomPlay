package user_id

import (
	"github.com/google/uuid"
)

type UserId uuid.UUID

func (id *UserId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}

func (id UserId) ToUuid() uuid.UUID {
	concreteId := uuid.UUID(id)
	return concreteId
}

func ParseUserId(s string) *UserId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := UserId(uuidResult)
	return &result
}

func NewUserId() UserId {
	return UserId(uuid.New())
}
