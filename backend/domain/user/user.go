package user

import (
	"github.com/google/uuid"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/shared"
)

const IdClaimContextKeyName = "userIdClaimContextKey"

type User struct {
	shared.AggregateRoot[shared.UserId]
	externalId string
	fullName   FullName
	role       *Role
	roomId     *shared.RoomId
	devices    []device.Device
}

func NewUser(externalId, name, surname string, deviceEntity device.Device) *User {
	user := &User{
		externalId: externalId,
		fullName:   *NewFullName(name, surname),
		role:       nil,
		roomId:     nil,
		devices:    []device.Device{deviceEntity},
	}
	user.SetId(shared.UserId(uuid.New()))
	return user
}
func (u *User) ChangeFullName(newFullName FullName) {
	u.fullName = newFullName
}
func (u *User) Devices() []device.Device {
	return u.devices
}
func (u *User) RoomId() *shared.RoomId {
	return u.roomId
}
func (u *User) FullName() FullName {
	return u.fullName
}
func (u *User) ExternalId() string {
	return u.externalId
}
func (u *User) Role() *Role {
	return u.role
}
func (u *User) LoginWithNewDevice(deviceType device.Type) shared.DeviceId {
	newDevice := device.NewDevice(deviceType)
	u.devices = append(u.devices, *newDevice)
	return newDevice.Id()
}
func (u *User) ReloginWithKnownDevice(deviceId shared.DeviceId) {
	u.getDeviceById(deviceId).RefreshDeviceState()
}
func (u *User) getDeviceById(deviceId shared.DeviceId) *device.Device {
	for i := range u.devices {
		if u.devices[i].Id() == deviceId {
			return &u.devices[i]
		}
	}
	return nil
}
func HydrateUser(
	id shared.UserId,
	externalId string,
	name string,
	surname string,
	role *Role,
	roomId *shared.RoomId,
	devices []device.Device,
) *User {
	user := &User{
		externalId: externalId,
		fullName:   *NewFullName(name, surname),
		role:       role,
		roomId:     roomId,
		devices:    devices,
	}
	user.SetId(id)
	return user
}
