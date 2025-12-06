package infrustructure_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

var userIds = []user.Id{
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
}
var songs = []room.Song{
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId1",
		"title1",
		"author1",
		349,
		userIds[0],
		time.Date(2001, 11, 22, 12, 00, 00, 00, time.UTC),
		nil,
		room.Played,
		88,
		false,
		false,
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId2",
		"title2",
		"author2",
		349,
		userIds[1],
		time.Date(2001, 10, 22, 12, 00, 00, 00, time.UTC),
		nil,
		room.Playing,
		8,
		true,
		false,
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId3",
		"title3",
		"author3",
		349,
		userIds[0],
		time.Date(2003, 10, 22, 12, 00, 00, 00, time.UTC),
		nil,
		room.Enqueued,
		0,
		false,
		true,
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId4",
		"title4",
		"author4",
		250,
		userIds[2],
		time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
		nil,
		room.Enqueued,
		5,
		false,
		false,
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId5",
		"title5",
		"author5",
		180,
		userIds[3],
		time.Date(2023, 5, 10, 12, 00, 00, 00, time.UTC),
		nil,
		room.Played,
		15,
		true,
		true,
	),
}
var rooms = []room.Room{
	*room.HydrateRoom(shared.RoomId(uuid.New()),
		"room1",
		"roompass1",
		"qrCode1",
		nil,
		time.Date(2001, 11, 12, 12, 00, 00, 00, time.UTC),
		uint32(time.Hour*30/time.Second),
		songs,
		userIds,
	),
	*room.HydrateRoom(shared.RoomId(uuid.New()),
		"room2",
		"roompass2",
		"qrCode2",
		nil,
		time.Date(2001, 11, 10, 12, 00, 00, 00, time.UTC),
		uint32(time.Hour*12/time.Second),
		songs,
		[]user.Id{userIds[0]},
	),
	*room.HydrateRoom(shared.RoomId(uuid.New()),
		"room3",
		"roompass3",
		"qrCode3",
		nil,
		time.Date(2022, 1, 1, 12, 00, 00, 00, time.UTC),
		uint32(time.Hour*24/time.Second),
		[]room.Song{songs[2], songs[3]},
		[]user.Id{userIds[1], userIds[2]},
	),
}
var devices = []user.Device{
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device1",
		user.Mobile,
		false,
		user.Offline,
		time.Date(2001, 12, 22, 12, 00, 00, 00, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device2",
		user.Desktop,
		true,
		user.Online,
		time.Date(2002, 12, 22, 12, 00, 00, 00, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device3",
		user.Mobile,
		true,
		user.Online,
		time.Date(2023, 1, 1, 12, 00, 00, 00, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device4",
		user.Desktop,
		false,
		user.Offline,
		time.Date(2023, 2, 2, 12, 00, 00, 00, time.UTC),
	),
}
var users = []*user.User{
	user.HydrateUser(userIds[0],
		"externalId1",
		"name1",
		"surname1",
		nil,
		nil,
		devices,
		nil),
	user.HydrateUser(userIds[1],
		"externalId2",
		"name2",
		"surname2",
		nil,
		nil,
		[]user.Device{devices[1]},
		nil),
	user.HydrateUser(userIds[2],
		"externalId3",
		"name3",
		"surname3",
		nil,
		nil,
		[]user.Device{devices[2]},
		nil),
	user.HydrateUser(userIds[3],
		"externalId4",
		"name4",
		"surname4",
		nil,
		nil,
		[]user.Device{devices[3]},
		nil),
	user.HydrateUser(userIds[4],
		"externalId5",
		"name5",
		"surname5",
		nil,
		nil,
		[]user.Device{devices[0]},
		nil),
}

var SeedData = struct {
	Room room.Room
	User user.User
}{}

type Seeder struct {
	Queryer contracts.IQueryer
}

func NewSeeder(queryer contracts.IQueryer) *Seeder {
	return &Seeder{
		Queryer: queryer,
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	for _, user := range users {
		if err := s.SeedUser(ctx, user); err != nil {
			return err
		}
		devices := user.Devices()
		for i := range devices {
			if err := s.SeedDevice(ctx, &devices[i], user.Id()); err != nil {
				return err
			}
		}
	}

	for i := range songs {
		if err := s.SeedSong(ctx, &songs[i]); err != nil {
			return err
		}
	}

	for i := range rooms {
		room := &rooms[i]
		if err := s.SeedRoom(ctx, room); err != nil {
			return err
		}
		songsInRoom := room.SongsList()
		for j := range songsInRoom {
			if err := s.SeedEnqueuedSong(ctx, &songsInRoom[j], room.Id()); err != nil {
				return err
			}
		}
		for _, memberId := range room.Members() {
			if err := s.SeedUserRoomData(ctx, room.Id(), memberId); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Seeder) SeedRoom(ctx context.Context, room *room.Room) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, room.Id(), room.Name(), []byte(room.Password()), []byte(room.QrCode()), room.CreatedAtUtc(), room.LifespanSeconds())
	if err != nil {
		return fmt.Errorf("failed to seed room: %w", err)
	}
	return nil
}

func (s *Seeder) SeedUser(ctx context.Context, user *user.User) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users (id, external_id, name, surname)
		VALUES ($1, $2, $3, $4)
	`, user.Id(), user.ExternalId(), user.FullName().Name(), user.FullName().Surname())
	if err != nil {
		return fmt.Errorf("failed to seed user: %w", err)
	}
	return nil
}

func (s *Seeder) SeedSong(ctx context.Context, song *room.Song) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO songs (id, external_id, title, author, length_seconds, album_cover_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, song.Id(), song.ExternalId(), song.Title(), song.Author(), song.LengthSeconds(), song.a())
	if err != nil {
		return fmt.Errorf("failed to seed song: %w", err)
	}
	return nil
}

func (s *Seeder) SeedEnqueuedSong(ctx context.Context, song *room.Song, roomID shared.RoomId) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO enqueued_songs (id, room_id, song_id, added_by, added_at_utc, state, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, uuid.New(), roomID, song.Id(), song.AddedBy(), song.StartedAtUtc(), song.State().String(), song.Votes())
	if err != nil {
		return fmt.Errorf("failed to seed enqueued song: %w", err)
	}
	return nil
}

func (s *Seeder) SeedDevice(ctx context.Context, device *user.Device, userID user.Id) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, device.Id(), device.FriendlyName(), device.IsHost(), device.DeviceType().String(), userID, device.State().String(), device.LastLoggedInUtc())
	if err != nil {
		return fmt.Errorf("failed to seed device: %w", err)
	}
	return nil
}

func (s *Seeder) SeedUserRoomData(ctx context.Context, roomID shared.RoomId, userID user.Id) error {
	role := user.Member
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users_room_data (room_id, user_id, role)
		VALUES ($1, $2, $3)
	`, roomID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to seed user_room_data: %w", err)
	}
	return nil
}
