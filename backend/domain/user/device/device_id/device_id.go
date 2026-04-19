package device_id

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

func (id DeviceId) ToUuid() uuid.UUID {
	concreteId := uuid.UUID(id)
	return concreteId
}

func ParseDeviceId(s string) *DeviceId {
	uuidResult, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	result := DeviceId(uuidResult)
	return &result
}

func NewDeviceId() DeviceId {
	return DeviceId(uuid.New())
}

func IdsEqual(id1, id2 *DeviceId) bool {
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
