package entities

import (
	"github.com/google/uuid"
)

type Type string

const (
	MOBILE   Type = "MOBILE"
	COMPUTER Type = "COMPUTER"
)

type DeviceId uuid.UUID

type Device struct {
	Entity[DeviceId]
	Fingerprint  string
	FriendlyName string
	Type         Type
	IsHost       bool
}
