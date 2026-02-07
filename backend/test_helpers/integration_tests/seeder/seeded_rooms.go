package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
)

var (
	roomIds = []room_id.RoomId{
		room_id.NewRoomId(),
		room_id.NewRoomId(),
		room_id.NewRoomId(),
	}

	rooms = []room.Room{
		*room.HydrateRoom(roomIds[0],
			"room1",
			"roompasss1",
			"qrCode1",
			nil,
			time.Date(2001, 11, 12, 12, 0o0, 0o0, 0o0, time.UTC),
			uint32(time.Hour*30/time.Second),
			enqueuedSongs,
			[]user_id.UserId{userIds[3]},
		),
		*room.HydrateRoom(roomIds[1],
			"room2",
			"roompasss2",
			"qrCode2",
			nil,
			time.Date(2001, 11, 10, 12, 0o0, 0o0, 0o0, time.UTC),
			uint32(time.Hour*12/time.Second),
			enqueuedSongs,
			[]user_id.UserId{userIds[0]},
		),
		*room.HydrateRoom(roomIds[2],
			"room3",
			"roompasss3",
			"qrCode3",
			nil,
			time.Date(2022, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
			uint32(time.Hour*24/time.Second),
			[]enqueued_song.EnqueuedSong{enqueuedSongs[2], enqueuedSongs[3]},
			[]user_id.UserId{userIds[1], userIds[2]},
		),
	}
)

func (s *Seeder) seedRoom(ctx context.Context,
	room *room.Room,
) error {
	configuration := mock_configuration.MockConfiguration{}
	encrypter := encryper.NewEncrypter(&configuration)
	hashedPassword, _ := encrypter.HashAndSalt(room.Password())
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO rooms 
		(
			id, 
			name,
			password,
			qr_code_hash,
			boost_cooldown_seconds,
			created_at_utc,
			lifespan_seconds
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7
		)
		`, room.Id().ToUuid(),
		room.Name(),
		[]byte(hashedPassword),
		[]byte(room.QrCode()),
		room.BoostCooldownSeconds(),
		room.CreatedAtUtc(),
		room.LifespanSeconds(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed room: %w", err)
	}

	return nil
}
