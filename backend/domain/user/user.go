package user

import (
	"github.com/google/uuid"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/shared"
)

const IdClaimContextKeyName = "userIdClaim"

type User struct {
	shared.AggregateRoot[shared.UserId]
	externalId string
	name       string
	surname    string
	role       Role
	roomId     *shared.RoomId
	devices    []device.Device
}

func FirstLogin(externalId, name, surname string, deviceEntity device.Device) *User {
	user := &User{
		externalId: externalId,
		name:       name,
		surname:    surname,
		role:       Guest,
		roomId:     nil,
		devices:    []device.Device{deviceEntity},
	}
	user.SetId(shared.UserId(uuid.New()))
	return user
}
func (u User) Devices() []device.Device {
	return u.devices
}

func (u User) RoomId() *shared.RoomId {
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
func (u User) Role() Role {
	return u.role
}
