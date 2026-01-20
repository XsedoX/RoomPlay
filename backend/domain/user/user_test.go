package user

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetMostRecentDevice(t *testing.T) {
	offlineDeviceState := Offline
	onlineDeviceState := Online
	deviceTypeMobile := Mobile
	deviceTypeDesktop := Desktop
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	device1 := HydrateDevice(DeviceId(uuid.New()),
		faker.Name(),
		deviceTypeDesktop,
		false,
		offlineDeviceState,
		time1)

	time2 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	device2 := HydrateDevice(DeviceId(uuid.New()),
		faker.Name(),
		deviceTypeMobile,
		true,
		onlineDeviceState,
		time2)

	role := Host
	boostUsedAtUtc := time.Now().UTC()
	roomId := shared.RoomId(uuid.New())
	user := HydrateUser(Id(uuid.New()),
		faker.UUIDDigit(),
		faker.Name(),
		faker.LastName(),
		&role,
		&roomId,
		[]Device{*device1, *device2},
		&boostUsedAtUtc)

	mostRecentDevice := user.GetMostRecentDevice()

	require.Equal(t, device2, mostRecentDevice, "The most recent device should be device2")
}

func TestLoginWithNewDevice(t *testing.T) {
	offlineDeviceState := Offline
	deviceTypeMobile := Mobile
	deviceTypeDesktop := Desktop
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	device1 := HydrateDevice(DeviceId(uuid.New()),
		faker.Name(),
		deviceTypeDesktop,
		false,
		offlineDeviceState,
		time1)

	role := Host
	boostUsedAtUtc := time.Now().UTC()
	roomId := shared.RoomId(uuid.New())
	user := HydrateUser(Id(uuid.New()),
		faker.UUIDDigit(),
		faker.Name(),
		faker.LastName(),
		&role,
		&roomId,
		[]Device{*device1},
		&boostUsedAtUtc)

	user.LoginWithNewDevice(deviceTypeMobile)

	require.Equal(t, 2, len(user.Devices()), "User should have 2 devices after logging in with a new device")
}

func TestNewUSer(t *testing.T) {
	externalID := faker.UUIDDigit()
	name := faker.Name()
	surname := faker.LastName()
	deviceType := Mobile
	newUser := NewUser(externalID, name, surname, deviceType)
	require.Equal(t, externalID, newUser.ExternalId(), "External ID should match")
	require.Equal(t, name+" "+surname, newUser.FullName().String(), "Full name should match")
	require.Nil(t, newUser.Role(), "Role should be nil for new user")
	require.Nil(t, newUser.RoomId(), "Room ID should be nil for new user")
	require.Nil(t, newUser.BoostUsedAtUtc(), "BoostUsedAtUtc should be nil for new user")
	require.Equal(t, 1, len(newUser.Devices()), "New user should have one device")
	require.Equal(t, deviceType, newUser.Devices()[0].DeviceType(), "Device type should match")
}
