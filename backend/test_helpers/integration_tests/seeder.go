package integration_tests

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication"
	othermocks "github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

var userIds = []user.Id{
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
	user.Id(uuid.New()),
}

var songsStartedAt = []time.Time{
	time.Date(2001, 11, 22, 12, 5, 0o0, 0o0, time.UTC),
	time.Date(2022, 12, 1, 15, 30, 0o0, 0o0, time.UTC),
	time.Date(2023, 6, 15, 18, 45, 0o0, 0o0, time.UTC),
}

var songs = []room.Song{
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId1",
		"title1",
		"author1",
		349,
		userIds[0],
		time.Date(2001, 11, 22, 12, 0o0, 0o0, 0o0, time.UTC),
		&songsStartedAt[0],
		room.Played,
		88,
		false,
		false,
		faker.URL(),
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId2",
		"title2",
		"author2",
		349,
		userIds[1],
		time.Date(2001, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
		&songsStartedAt[1],
		room.Playing,
		8,
		true,
		false,
		faker.URL(),
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId3",
		"title3",
		"author3",
		349,
		userIds[0],
		time.Date(2003, 10, 22, 12, 0o0, 0o0, 0o0, time.UTC),
		nil,
		room.Enqueued,
		0,
		false,
		true,
		faker.URL(),
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId4",
		"title4",
		"author4",
		250,
		userIds[2],
		time.Date(2024, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
		nil,
		room.Enqueued,
		5,
		false,
		false,
		faker.URL(),
	),
	*room.HydrateSong(room.SongId(uuid.New()),
		"songExternalId5",
		"title5",
		"author5",
		180,
		userIds[3],
		time.Date(2023, 5, 10, 12, 0o0, 0o0, 0o0, time.UTC),
		&songsStartedAt[2],
		room.Played,
		15,
		true,
		true,
		faker.URL(),
	),
}

var (
	roomId0 = shared.RoomId(uuid.New())
	roomId1 = shared.RoomId(uuid.New())
	roomId2 = shared.RoomId(uuid.New())
)

var rooms = []room.Room{
	*room.HydrateRoom(roomId0,
		"room1",
		"roompasss1",
		"qrCode1",
		nil,
		time.Date(2001, 11, 12, 12, 0o0, 0o0, 0o0, time.UTC),
		uint32(time.Hour*30/time.Second),
		songs,
		[]user.Id{userIds[3]},
	),
	*room.HydrateRoom(roomId1,
		"room2",
		"roompasss2",
		"qrCode2",
		nil,
		time.Date(2001, 11, 10, 12, 0o0, 0o0, 0o0, time.UTC),
		uint32(time.Hour*12/time.Second),
		songs,
		[]user.Id{userIds[0]},
	),
	*room.HydrateRoom(roomId2,
		"room3",
		"roompasss3",
		"qrCode3",
		nil,
		time.Date(2022, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
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
		time.Date(2001, 12, 22, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device2",
		user.Desktop,
		true,
		user.Online,
		time.Date(2002, 12, 22, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device3",
		user.Mobile,
		true,
		user.Online,
		time.Date(2023, 1, 1, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device4",
		user.Desktop,
		false,
		user.Offline,
		time.Date(2023, 2, 2, 12, 0o0, 0o0, 0o0, time.UTC),
	),
	*user.HydrateDevice(user.DeviceId(uuid.New()),
		"device5",
		user.Desktop,
		false,
		user.Online,
		time.Date(2025, 2, 2, 12, 0o0, 0o0, 0o0, time.UTC),
	),
}

var (
	user1Role = user.Host
	user2Role = user.Member
	user3Role = user.Member
	user4Role = user.Member
	user5Role = user.Member
	users     = []user.User{
		*user.HydrateUser(userIds[0],
			"externalId1",
			"name1",
			"surname1",
			&user1Role,
			&roomId1,
			[]user.Device{devices[4]},
			nil),
		*user.HydrateUser(userIds[1],
			"externalId2",
			"name2",
			"surname2",
			&user2Role,
			&roomId2,
			[]user.Device{devices[1]},
			nil),
		*user.HydrateUser(userIds[2],
			"externalId3",
			"name3",
			"surname3",
			&user3Role,
			&roomId2,
			[]user.Device{devices[2]},
			nil),
		*user.HydrateUser(userIds[3],
			"externalId4",
			"name4",
			"surname4",
			&user4Role,
			&roomId2,
			[]user.Device{devices[3]},
			nil),
		*user.HydrateUser(userIds[4],
			"externalId5",
			"name5",
			"surname5",
			&user5Role,
			nil,
			[]user.Device{devices[0]},
			nil),
	}
)

var refreshTokens = []credentials.RefreshToken{
	*credentials.HydrateRefreshToken(users[0].Id(),
		users[0].Devices()[0].Id(),
		"refreshTokenValue1",
		time.Now().AddDate(1, 0, 0),
		time.Date(2023, 12, 1, 12, 0o0, 0o0, 0o0, time.UTC)),
}

var SeedData = struct {
	Rooms                    []room.Room
	Users                    []user.User
	Songs                    []room.Song
	Devices                  []user.Device
	LoggedInUserRefreshToken credentials.RefreshToken
}{
	Rooms:                    rooms,
	Songs:                    songs,
	Users:                    users,
	Devices:                  devices,
	LoggedInUserRefreshToken: refreshTokens[0],
}

type Seeder struct {
	Queryer contracts.IQueryer
}

func NewSeeder(queryer contracts.IQueryer) *Seeder {
	return &Seeder{
		Queryer: queryer,
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	for _, user1 := range users {
		if err := s.SeedUser(ctx, &user1); err != nil {
			return err
		}
		devices1 := user1.Devices()
		for i := range devices1 {
			if err := s.SeedDevice(ctx, &devices1[i], user1.Id()); err != nil {
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
		room1 := &rooms[i]
		if err := s.SeedRoom(ctx, room1); err != nil {
			return err
		}
		songsInRoom := room1.SongsList()
		for j := range songsInRoom {
			if err := s.SeedEnqueuedSong(ctx, &songsInRoom[j], room1.Id()); err != nil {
				return err
			}
		}
		for _, memberId := range room1.Members() {
			var userRole *user.UserRole
			for _, u := range users {
				if u.Id() == memberId {
					userRole = u.Role()
					break
				}
			}
			if err := s.SeedUserRoomData(ctx, room1.Id(), memberId, userRole); err != nil {
				return err
			}
		}
	}
	for _, usersRefreshToken := range refreshTokens {
		if err := s.SeedRefreshToken(ctx, usersRefreshToken); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) SeedRoom(ctx context.Context, room *room.Room) error {
	configuration := othermocks.MockConfiguration{}
	encrypter := authentication.NewEncrypter(&configuration)
	hashedPassword, _ := encrypter.HashAndSalt(room.Password())
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO rooms (id, name, password, qr_code_hash, created_at_utc, lifespan_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, room.Id(), room.Name(), []byte(hashedPassword), []byte(room.QrCode()), room.CreatedAtUtc(), room.LifespanSeconds())
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
	`, song.Id(), song.ExternalId(), song.Title(), song.Author(), song.LengthSeconds(), song.AlbumCoverUrl())
	if err != nil {
		return fmt.Errorf("failed to seed song: %w", err)
	}
	return nil
}

func (s *Seeder) SeedEnqueuedSong(ctx context.Context, song *room.Song, roomID shared.RoomId) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO enqueued_songs (id, room_id, song_id, added_by, added_at_utc, started_at_utc, state, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, uuid.New(), roomID, song.Id(), song.AddedBy(), song.AddedAtUtc(), song.StartedAtUtc(), song.State().String(), song.Votes())
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

func (s *Seeder) SeedUserRoomData(ctx context.Context, roomID shared.RoomId, userID user.Id, role *user.UserRole) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users_room_data (room_id, user_id, role)
		VALUES ($1, $2, $3)
	`, roomID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to seed user_room_data: %w", err)
	}
	return nil
}

func (s *Seeder) SeedRefreshToken(ctx context.Context, refreshToken credentials.RefreshToken) error {
	configuration := othermocks.MockConfiguration{}
	encrypter := authentication.NewEncrypter(&configuration)
	hashedRefreshToken := encrypter.Hash(refreshToken.RefreshToken())
	_, err := s.Queryer.ExecContext(ctx, `
INSERT INTO users_refresh_tokens (user_id, device_id, refresh_token, expires_at_utc, issued_at_utc)
VALUES ($1, $2, $3, $4, $5)
`, SeedData.LoggedInUserRefreshToken.Id(),
		SeedData.LoggedInUserRefreshToken.DeviceId(),
		hashedRefreshToken,
		SeedData.LoggedInUserRefreshToken.ExpiresAtUtc(),
		SeedData.LoggedInUserRefreshToken.IssuedAtUtc())
	if err != nil {
		return fmt.Errorf("failed to seed refresh token: %w", err)
	}
	return nil
}
