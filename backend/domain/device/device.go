package device

import (
	"github.com/google/uuid"
	"xsedox.com/main/domain/shared"
)

type Device struct {
	shared.Entity[shared.DeviceId]
	friendlyName string
	deviceType   Type
	isHost       bool
	state        State
}

func (d Device) FriendlyName() string {
	return d.friendlyName
}

func (d Device) DeviceType() Type {
	return d.deviceType
}

func (d Device) IsHost() bool {
	return d.isHost
}

func (d Device) State() State {
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
		friendlyName: friendlyName,
		deviceType:   deviceType,
		state:        Online,
		isHost:       false,
	}
	device.SetId(shared.DeviceId(uuid.New()))
	return device
}
