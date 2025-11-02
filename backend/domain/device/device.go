package device

import (
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/shared"
)

type Device struct {
	shared.Entity[shared.DeviceId]
	friendlyName    string
	deviceType      Type
	isHost          bool
	state           State
	lastLoggedInUtc time.Time
}

func (d *Device) LastLoggedInUtc() time.Time {
	return d.lastLoggedInUtc
}

func (d *Device) FriendlyName() string {
	return d.friendlyName
}

func (d *Device) DeviceType() Type {
	return d.deviceType
}

func (d *Device) IsHost() bool {
	return d.isHost
}

func (d *Device) State() State {
	return d.state
}

func NewDevice(deviceType Type) *Device {
	var friendlyName string
	switch deviceType {
	case Mobile:
		friendlyName = "My lovely mobile device"
	case Computer:
		friendlyName = "My lovely computer device"
	}
	device := &Device{
		friendlyName:    friendlyName,
		deviceType:      deviceType,
		state:           Online,
		isHost:          false,
		lastLoggedInUtc: time.Now().UTC(),
	}
	device.SetId(shared.DeviceId(uuid.New()))
	return device
}
func (d *Device) RefreshDeviceState() {
	d.state = Online
	d.lastLoggedInUtc = time.Now().UTC()
}
func (d *Device) ChangeDeviceFriendlyName(friendlyName string) {
	d.friendlyName = friendlyName
}
func (d *Device) ChangeDeviceType(deviceType Type) {
	d.deviceType = deviceType
}
func HydrateDevice(
	id shared.DeviceId,
	friendlyName string,
	deviceType Type,
	isHost bool,
	state State,
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
