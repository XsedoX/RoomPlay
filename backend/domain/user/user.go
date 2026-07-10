package user

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/full_name"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
)

const IdClaimContextKeyName = "userIdClaimContextKey"

type User struct {
	shared.AggregateRoot[user_id.UserId]
	fullName       full_name.FullName
	role           *user_role.UserRole
	roomId         *room_id.RoomId
	devices        []device.Device
	boostUsedAtUtc *time.Time
}

func NewUser(name, surname string, deviceType device_type.DeviceType) *User {
	deviceEntity := device.NewDevice(deviceType)
	user := &User{
		fullName:       full_name.NewFullName(name, surname),
		role:           nil,
		roomId:         nil,
		devices:        []device.Device{*deviceEntity},
		boostUsedAtUtc: nil,
	}
	user.SetId(user_id.NewUserId())
	return user
}

func (u *User) GetMostRecentDevice() *device.Device {
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

func (u *User) ChangeFullName(newFullName full_name.FullName) {
	u.fullName = newFullName
}

func (u *User) Devices() []device.Device {
	return u.devices
}

func (u *User) RoomId() *room_id.RoomId {
	return u.roomId
}

func (u *User) FullName() full_name.FullName {
	return u.fullName
}

func (u *User) Role() *user_role.UserRole {
	return u.role
}

func (u *User) LoginWithNewDevice(deviceType device_type.DeviceType) device_id.DeviceId {
	newDevice := device.NewDevice(deviceType)
	u.devices = append(u.devices, *newDevice)
	return newDevice.Id()
}

func (u *User) CheckDeviceOwnership(deviceId device_id.DeviceId) bool {
	usersDevice := u.getDeviceById(deviceId)
	return usersDevice != nil
}

func (u *User) ReloginWithKnownDevice(deviceId device_id.DeviceId) error {
	usersDevice := u.getDeviceById(deviceId)
	if usersDevice == nil {
		return domain_errors.NewUserDeviceNotFoundError(u.Id(), deviceId)
	}
	usersDevice.RefreshDeviceState()
	return nil
}

func (u *User) getDeviceById(deviceId device_id.DeviceId) *device.Device {
	for i := range u.devices {
		usersDeviceId := u.devices[i].Id()
		if usersDeviceId == deviceId {
			return &u.devices[i]
		}
	}
	return nil
}

func HydrateUser(
	id user_id.UserId,
	name string,
	surname string,
	role *user_role.UserRole,
	roomId *room_id.RoomId,
	devices []device.Device,
	boostUsedAtUtc *time.Time,
) *User {
	user := &User{
		fullName:       full_name.NewFullName(name, surname),
		role:           role,
		roomId:         roomId,
		devices:        devices,
		boostUsedAtUtc: boostUsedAtUtc,
	}
	user.SetId(id)
	return user
}
