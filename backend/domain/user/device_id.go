package user

import (
	"github.com/google/uuid"
)

type DeviceId uuid.UUID

func (id *DeviceId) String() *string {
	if id == nil {
		return nil
	}
	uuId := uuid.UUID(*id).String()
	return &uuId
}
func (id *DeviceId) ToUuid() *uuid.UUID {
	if id == nil {
		return nil
	}
	concreteId := uuid.UUID(*id)
	return &concreteId
}
func ParseDeviceId(s string) *DeviceId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := DeviceId(uuidResult)
	return &result
}
