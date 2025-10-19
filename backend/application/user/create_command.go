package user

import "xsedox.com/main/domain/entities"

type CreateCommand struct {
	Name       string    `json:"name" validate:"required"`
	Device     DeviceDTO `json:"device" validate:"required,dive"`
	ExternalId string    `json:"external_id" validate:"required"`
	Surname    string    `json:"surname" validate:"required"`
}
type DeviceDTO struct {
	Fingerprint  string              `json:"fingerprint" validate:"required"`
	FriendlyName string              `json:"friendly_name" validate:"required"`
	DeviceType   entities.DeviceType `json:"device_type" validate:"required, device_type_validation"`
}
