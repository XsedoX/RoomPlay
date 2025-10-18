package entities

import (
	"github.com/google/uuid"
	"xsedox.com/domain/value_objects"
)

type Role string

const (
	HOST Role = "ROOM_CREATOR"
	USER Role = "USER"
)

type UserId uuid.UUID

type User struct {
	Entity[UserId]
	externalId string
	email      value_objects.Email
	role       Role
	room       *Room
	devices    []Device
}

func NewUser(externalId string, emailString string, role Role, deviceParam Device) (*User, error) {
	email, err := value_objects.NewEmail(emailString)
	if err != nil {
		return nil, err
	}
	user := &User{
		externalId: externalId,
		email:      *email,
		role:       role,
		room:       nil,
		devices:    []Device{deviceParam},
	}
	user.SetId(UserId(uuid.New()))
	return user, nil
}
