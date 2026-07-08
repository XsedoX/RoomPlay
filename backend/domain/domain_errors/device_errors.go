package domain_errors

import "fmt"

func NewDeviceFriendlyNameWrongLengthError(deviceNameMaxLength uint8) error {
	return &DomainError{
		Code:        "Device.FriendlyName.WrongLength",
		Description: fmt.Sprintf("The device friendly name exceeded %d characters.", deviceNameMaxLength),
	}
}
