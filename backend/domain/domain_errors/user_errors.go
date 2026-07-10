package domain_errors

import (
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

func NewUserDeviceNotFoundError(userId user_id.UserId, deviceId device_id.DeviceId) *DomainError {
	return &DomainError{
		Code:        "User.Device.NotFound",
		Description: "Device with id " + *deviceId.String() + " not found for user with id " + *userId.String(),
	}
}
