package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

var devices = []device.Device{
	*device.HydrateDevice(device_id.NewDeviceId(),
		"device1",
		device_type.Mobile,
		false,
		device_state.Offline,
		time.Date(2001, 12, 22, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*device.HydrateDevice(device_id.NewDeviceId(),
		"device2",
		device_type.Desktop,
		true,
		device_state.Online,
		time.Date(2002, 12, 22, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*device.HydrateDevice(device_id.NewDeviceId(),
		"device3",
		device_type.Mobile,
		true,
		device_state.Online,
		time.Date(2023, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*device.HydrateDevice(device_id.NewDeviceId(),
		"device4",
		device_type.Desktop,
		false,
		device_state.Offline,
		time.Date(2023, 2, 2, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*device.HydrateDevice(device_id.NewDeviceId(),
		"device5",
		device_type.Desktop,
		false,
		device_state.Online,
		time.Date(2025, 2, 2, 12, 0o0, 0o0, 0o0, time.UTC),
	),
}

func (s *Seeder) seedDevice(ctx context.Context, device *device.Device, userID user_id.UserId) error {
	_, err := s.Queryer.ExecContext(ctx,
		`
		INSERT INTO devices (
			id,
			friendly_name,
			is_host,
			type,
			user_id,
			state,
			last_logged_in_at_utc
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`, device.Id(),
		device.FriendlyName(),
		device.IsHost(),
		device.DeviceType().String(),
		userID,
		device.State().String(),
		device.LastLoggedInUtc(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed device: %w", err)
	}
	return nil
}
