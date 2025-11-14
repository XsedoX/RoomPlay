package user

import (
	"github.com/google/uuid"
)

type Id uuid.UUID

func (id *Id) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}
func (id *Id) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	concreteId := uuid.UUID(*id)
	return &concreteId
}
func ParseUserId(s string) *Id {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := Id(uuidResult)
	return &result
}
