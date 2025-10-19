package entities

import (
	"github.com/google/uuid"
)

type DeviceType string

const (
	MOBILE   DeviceType = "MOBILE"
	COMPUTER DeviceType = "COMPUTER"
)

type DeviceId uuid.UUID

type Device struct {
	Entity[DeviceId]
	fingerprint  string
	friendlyName string
	deviceType   DeviceType
	isHost       bool
}

func NewDevice(fingerPrint, friendlyName string, deviceType DeviceType) *Device {
	device := &Device{
		fingerprint:  fingerPrint,
		friendlyName: friendlyName,
		deviceType:   deviceType,
		isHost:       false,
	}
	device.SetId(DeviceId(uuid.New()))
	return device
}
