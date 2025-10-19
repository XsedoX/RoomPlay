package entities

import (
	"github.com/google/uuid"
)

type UserRole string

const (
	HOST UserRole = "HOST"
	USER UserRole = "USER"
)

type UserId uuid.UUID

type User struct {
	Entity[UserId]
	externalId string
	name       string
	surname    string
	roomId     *RoomId
	devices    []Device
}

func FirstLogin(externalId, name, surname string, device Device) *User {
	user := &User{
		externalId: externalId,
		name:       name,
		surname:    surname,
		roomId:     nil,
		devices:    []Device{device},
	}
	user.SetId(UserId(uuid.New()))
	return user
}
func (u User) Devices() []Device {
	return u.devices
}

func (u User) RoomId() *RoomId {
	return u.roomId
}

func (u User) Surname() string {
	return u.surname
}

func (u User) Name() string {
	return u.name
}

func (u User) ExternalId() string {
	return u.externalId
}
