package device

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestNewDesktopDevice(t *testing.T) {
	deviceType := device_type.Desktop

	newDevice := NewDevice(deviceType)

	require.Equal(t, deviceType, newDevice.DeviceType())
	require.Equal(t, device_state.Online, newDevice.State())
	require.False(t, newDevice.IsHost())
	require.Equal(t, "My lovely computer device", newDevice.FriendlyName())
	require.WithinDuration(t, time.Now().UTC(), newDevice.LastLoggedInUtc(), time.Second)
}

func TestNewMobileDevice(t *testing.T) {
	deviceType := device_type.Mobile

	newDevice := NewDevice(deviceType)

	require.Equal(t, deviceType, newDevice.DeviceType())
	require.Equal(t, device_state.Online, newDevice.State())
	require.False(t, newDevice.IsHost())
	require.Equal(t, "My lovely mobile device", newDevice.FriendlyName())
	require.WithinDuration(t, time.Now().UTC(), newDevice.LastLoggedInUtc(), time.Second)
}

func TestRefreshDeviceState(t *testing.T) {
	deviceType := device_type.Mobile
	friendlyName := faker.Word()
	isHost := false
	state := device_state.Offline
	deviceId := device_id.NewDeviceId()
	lastLoggedInUtc := time.Now().UTC().Add(-10 * time.Minute)
	newDevice := HydrateDevice(deviceId, friendlyName, deviceType, isHost, state, lastLoggedInUtc)

	newDevice.RefreshDeviceState()

	require.WithinDuration(t, time.Now().UTC(), newDevice.LastLoggedInUtc(), time.Second)
	require.Equal(t, device_state.Online, newDevice.State())
}

func TestChangeDeviceFriendlyName(t *testing.T) {
	deviceType := device_type.Desktop

	newDevice := NewDevice(deviceType)

	// Valid name
	validName := "ValidName"
	err := newDevice.ChangeDeviceFriendlyName(validName)
	require.NoError(t, err, "expected no error for valid name, got %v", err)
	require.Equal(t, validName, newDevice.FriendlyName())

	// Too short
	shortName := ""
	err = newDevice.ChangeDeviceFriendlyName(shortName)
	require.Error(t, err, "expected error for too short name, got nil")

	// Too long
	longName := faker.Password(options.WithRandomStringLength(31))
	err = newDevice.ChangeDeviceFriendlyName(longName)
	require.Error(t, err, "expected error for too long name, got nil")
}
