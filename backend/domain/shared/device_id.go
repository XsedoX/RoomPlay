package shared

import "github.com/google/uuid"

type DeviceId uuid.UUID

func (id DeviceId) String() string {
	return uuid.UUID(id).String()
}

// ParseDeviceId returns nil if something went wrong
func ParseDeviceId(deviceId string) *DeviceId {
	result, err := uuid.Parse(deviceId)
	if err != nil {
		return nil
	}
	deviceIdRes := DeviceId(result)
	return &deviceIdRes
}
