package device

import (
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
)

const (
	DeviceNameMaxLength = 30
	DeviceNameMinLength = 1
)

type Device struct {
	shared.Entity[device_id.DeviceId]
	friendlyName    string
	deviceType      device_type.DeviceType
	isHost          bool
	state           device_state.DeviceState
	lastLoggedInUtc time.Time
}

func (d *Device) LastLoggedInUtc() time.Time {
	return d.lastLoggedInUtc
}

func (d *Device) FriendlyName() string {
	return d.friendlyName
}

func (d *Device) DeviceType() device_type.DeviceType {
	return d.deviceType
}

func (d *Device) IsHost() bool {
	return d.isHost
}

func (d *Device) State() device_state.DeviceState {
	return d.state
}

func NewDevice(deviceType device_type.DeviceType) *Device {
	var friendlyName string
	switch deviceType {
	case device_type.Mobile:
		friendlyName = "My lovely mobile device"
	case device_type.Desktop:
		friendlyName = "My lovely computer device"
	}
	device := &Device{
		friendlyName:    friendlyName,
		deviceType:      deviceType,
		state:           device_state.Online,
		isHost:          false,
		lastLoggedInUtc: time.Now().UTC(),
	}
	device.SetId(device_id.NewDeviceId())
	return device
}

func (d *Device) RefreshDeviceState() {
	d.state = device_state.Online
	d.lastLoggedInUtc = time.Now().UTC()
}

func (d *Device) ChangeDeviceFriendlyName(friendlyName string) error {
	if len(friendlyName) > DeviceNameMaxLength {
		return validation_domain_error.NewValidationDomainError("Device.TooLong.FriendlyName", fmt.Sprintf("The device friendly name exceeded %d characters.", DeviceNameMaxLength))
	}
	if len(friendlyName) < DeviceNameMinLength {
		return validation_domain_error.NewValidationDomainError("Device.TooShort.FriendlyName", fmt.Sprintf("The device friendly name was shorter than %d characters.", DeviceNameMinLength))
	}
	d.friendlyName = friendlyName
	return nil
}

func (d *Device) ChangeDeviceType(deviceType device_type.DeviceType) {
	d.deviceType = deviceType
}

func HydrateDevice(
	id device_id.DeviceId,
	friendlyName string,
	deviceType device_type.DeviceType,
	isHost bool,
	state device_state.DeviceState,
	lastLoggedInUtc time.Time,
) *Device {
	device := &Device{
		friendlyName:    friendlyName,
		deviceType:      deviceType,
		isHost:          isHost,
		state:           state,
		lastLoggedInUtc: lastLoggedInUtc,
	}
	device.SetId(id)
	return device
}
