package shared

import "github.com/google/uuid"

type DeviceId uuid.UUID

func (id DeviceId) String() string {
	return uuid.UUID(id).String()
}
func ParseDeviceId(deviceId string) (*DeviceId, error) {
	result, err := uuid.Parse(deviceId)
	if err != nil {
		return nil, err
	}
	deviceIdRes := DeviceId(result)
	return &deviceIdRes, nil
}
