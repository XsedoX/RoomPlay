package domain_errors

import (
	"fmt"
)

func NewRoomNameIncorrectError(roomNameMaxLength, roomNameMinLength uint8) error {
	return &DomainError{
		Code: "Room.Name.WrongLength",
		Description: fmt.Sprintf(
			"The room name has to be between %d and %d characters.",
			roomNameMinLength,
			roomNameMaxLength,
		),
	}
}

func NewRoomPasswordIncorrectError(roomPasswordMaxLength, roomPasswordMinLength uint8) error {
	return &DomainError{
		Code: "Room.Password.WrongLength",
		Description: fmt.Sprintf(
			"The room password has to be between %d and %d characters.",
			roomPasswordMinLength,
			roomPasswordMaxLength,
		),
	}
}

func NewRoomQrCodeEmptyError() error {
	return &DomainError{
		Code:        "Room.QrCode.Empty",
		Description: "The room QR code cannot be empty.",
	}
}

func NewRoomHashSaltError() error {
	return &DomainError{
		Code:        "Room.HashSalt.Error",
		Description: "An error occurred while hashing and salting the room password.",
	}
}
