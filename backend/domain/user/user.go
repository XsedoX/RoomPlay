package user

import (
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/shared"
)

const IdClaimContextKeyName = "userIdClaimContextKey"

type User struct {
	shared.AggregateRoot[Id]
	externalId     string
	fullName       FullName
	role           *UserRole
	roomId         *shared.RoomId
	devices        []Device
	boostUsedAtUtc *time.Time
}

func NewUser(externalId, name, surname string, deviceType DeviceType) *User {
	deviceEntity := NewDevice(deviceType)
	user := &User{
		externalId:     externalId,
		fullName:       NewFullName(name, surname),
		role:           nil,
		roomId:         nil,
		devices:        []Device{*deviceEntity},
		boostUsedAtUtc: nil,
	}
	user.SetId(Id(uuid.New()))
	return user
}
func (u *User) GetMostRecentDevice() *Device {
	mostRecent := &u.devices[0]
	for i := range u.devices {
		if u.devices[i].LastLoggedInUtc().After(mostRecent.LastLoggedInUtc()) {
			mostRecent = &u.devices[i]
		}
	}
	return mostRecent
}
func (u *User) BoostUsedAtUtc() *time.Time {
	return u.boostUsedAtUtc
}
func (u *User) ChangeFullName(newFullName FullName) {
	u.fullName = newFullName
}
func (u *User) Devices() []Device {
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
func (u *User) Role() *UserRole {
	return u.role
}
func (u *User) LoginWithNewDevice(deviceType DeviceType) DeviceId {
	newDevice := NewDevice(deviceType)
	u.devices = append(u.devices, *newDevice)
	return newDevice.Id()
}
func (u *User) ReloginWithKnownDevice(deviceId DeviceId) {
	u.getDeviceById(deviceId).RefreshDeviceState()
}
func (u *User) getDeviceById(deviceId DeviceId) *Device {
	for i := range u.devices {
		if u.devices[i].Id() == deviceId {
			return &u.devices[i]
		}
	}
	return nil
}
func HydrateUser(
	id Id,
	externalId string,
	name string,
	surname string,
	role *UserRole,
	roomId *shared.RoomId,
	devices []Device,
	boostUsedAtUtc *time.Time,
) *User {
	user := &User{
		externalId:     externalId,
		fullName:       NewFullName(name, surname),
		role:           role,
		roomId:         roomId,
		devices:        devices,
		boostUsedAtUtc: boostUsedAtUtc,
	}
	user.SetId(id)
	return user
}
