package device_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/google/uuid"
)

type DeviceDao struct {
	Id                uuid.UUID `gorm:"primaryKey;type:uuid;" db:"id"`
	FriendlyName      string    `gorm:"friendly_name" db:"friendly_name"`
	IsHost            bool      `gorm:"is_host" db:"is_host"`
	Type              string    `gorm:"type:device_type" db:"type"`
	UserId            uuid.UUID `gorm:"user_id;foreignKey:UserId;references:Id" db:"user_id"`
	State             string    `gorm:"type:device_state;default:'online'"`
	LastLoggedInAtUtc time.Time `gorm:"last_logged_in_at_utc" db:"last_logged_in_at_utc"`
}

func (DeviceDao) TableName() string {
	return "devices"
}

func (d DeviceDao) ToDomain() *device.Device {
	return device.HydrateDevice(
		device_id.DeviceId(d.Id),
		d.FriendlyName,
		*device_type.ParseDeviceType(&d.Type),
		d.IsHost,
		*device_state.ParseDeviceState(&d.State),
		d.LastLoggedInAtUtc,
	)
}

// FromDomain takes the owner ID because Device does not store it.
func (d *DeviceDao) FromDomain(value *device.Device, userId user_id.UserId) {
	*d = DeviceDao{
		Id:                value.Id().ToUuid(),
		FriendlyName:      value.FriendlyName(),
		IsHost:            value.IsHost(),
		Type:              value.DeviceType().String(),
		UserId:            userId.ToUuid(),
		State:             value.State().String(),
		LastLoggedInAtUtc: value.LastLoggedInUtc(),
	}
}
