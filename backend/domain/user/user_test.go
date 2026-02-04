package user

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestGetMostRecentDevice(t *testing.T) {
	offlineDeviceState := device_state.Offline
	onlineDeviceState := device_state.Online
	deviceTypeMobile := device_type.Mobile
	deviceTypeDesktop := device_type.Desktop
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	device1 := device.HydrateDevice(device_id.NewDeviceId(),
		faker.Name(),
		deviceTypeDesktop,
		false,
		offlineDeviceState,
		time1)

	time2 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	device2 := device.HydrateDevice(device_id.NewDeviceId(),
		faker.Name(),
		deviceTypeMobile,
		true,
		onlineDeviceState,
		time2)

	role := user_role.Host
	boostUsedAtUtc := time.Now().UTC()
	roomId := room_id.NewRoomId()
	user := HydrateUser(user_id.NewUserId(),
		faker.Name(),
		faker.LastName(),
		&role,
		&roomId,
		[]device.Device{*device1, *device2},
		&boostUsedAtUtc)

	mostRecentDevice := user.GetMostRecentDevice()

	require.Equal(t, device2, mostRecentDevice, "The most recent device should be device2")
}

func TestLoginWithNewDevice(t *testing.T) {
	offlineDeviceState := device_state.Offline
	deviceTypeMobile := device_type.Mobile
	deviceTypeDesktop := device_type.Desktop
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	device1 := device.HydrateDevice(device_id.NewDeviceId(),
		faker.Name(),
		deviceTypeDesktop,
		false,
		offlineDeviceState,
		time1)

	role := user_role.Host
	boostUsedAtUtc := time.Now().UTC()
	roomId := room_id.NewRoomId()
	user := HydrateUser(user_id.NewUserId(),
		faker.Name(),
		faker.LastName(),
		&role,
		&roomId,
		[]device.Device{*device1},
		&boostUsedAtUtc)

	user.LoginWithNewDevice(deviceTypeMobile)

	require.Equal(t, 2, len(user.Devices()), "User should have 2 devices after logging in with a new device")
}

func TestNewUser(t *testing.T) {
	name := faker.Name()
	surname := faker.LastName()
	deviceType := device_type.Mobile
	newUser := NewUser(name, surname, deviceType)
	require.Equal(t, name+" "+surname, newUser.FullName().String(), "Full name should match")
	require.Nil(t, newUser.Role(), "Role should be nil for new user")
	require.Nil(t, newUser.RoomId(), "Room ID should be nil for new user")
	require.Nil(t, newUser.BoostUsedAtUtc(), "BoostUsedAtUtc should be nil for new user")
	require.Equal(t, 1, len(newUser.Devices()), "New user should have one device")
	require.Equal(t, deviceType, newUser.Devices()[0].DeviceType(), "Device type should match")
}
