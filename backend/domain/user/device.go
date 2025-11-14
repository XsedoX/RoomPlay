package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/domainErrors"
	"xsedox.com/main/domain/shared"
)

const (
	DeviceNameMaxLength = 30
	DeviceNameMinLength = 1
)

type Device struct {
	shared.Entity[DeviceId]
	friendlyName    string
	deviceType      DeviceType
	isHost          bool
	state           DeviceState
	lastLoggedInUtc time.Time
}

func (d *Device) LastLoggedInUtc() time.Time {
	return d.lastLoggedInUtc
}

func (d *Device) FriendlyName() string {
	return d.friendlyName
}

func (d *Device) DeviceType() DeviceType {
	return d.deviceType
}

func (d *Device) IsHost() bool {
	return d.isHost
}

func (d *Device) State() DeviceState {
	return d.state
}

func NewDevice(deviceType DeviceType) *Device {
	var friendlyName string
	switch deviceType {
	case Mobile:
		friendlyName = "My lovely mobile device"
	case Desktop:
		friendlyName = "My lovely computer device"
	}
	device := &Device{
		friendlyName:    friendlyName,
		deviceType:      deviceType,
		state:           Online,
		isHost:          false,
		lastLoggedInUtc: time.Now().UTC(),
	}
	device.SetId(DeviceId(uuid.New()))
	return device
}
func (d *Device) RefreshDeviceState() {
	d.state = Online
	d.lastLoggedInUtc = time.Now().UTC()
}
func (d *Device) ChangeDeviceFriendlyName(friendlyName string) error {
	if len(friendlyName) > DeviceNameMaxLength {
		return domainErrors.NewValidationDomainError("Device.TooLong.FriendlyName", fmt.Sprintf("The device friendly name exceeded %d characters.", DeviceNameMaxLength))
	}
	if len(friendlyName) < DeviceNameMinLength {
		return domainErrors.NewValidationDomainError("Device.TooShort.FriendlyName", fmt.Sprintf("The device friendly name was shorter than %d characters.", DeviceNameMinLength))
	}
	d.friendlyName = friendlyName
	return nil
}
func (d *Device) ChangeDeviceType(deviceType DeviceType) {
	d.deviceType = deviceType
}
func HydrateDevice(
	id DeviceId,
	friendlyName string,
	deviceType DeviceType,
	isHost bool,
	state DeviceState,
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
