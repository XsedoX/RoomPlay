package user

import (
	"xsedox.com/main/domain/entities"
)

type LoginCommand struct {
	Name       string    `json:"name" validate:"required"`
	Device     DeviceDto `json:"device" validate:"required,dive"`
	ExternalId string    `json:"external_id" validate:"required"`
	Surname    string    `json:"surname" validate:"required"`
}
type DeviceDto struct {
	Fingerprint string              `json:"fingerprint" validate:"required"`
	DeviceType  entities.DeviceType `json:"device_type" validate:"required, device_type_validation"`
}
