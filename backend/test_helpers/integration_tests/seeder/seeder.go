package seeder

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
)

var SeedData = struct {
	Rooms                    []room.Room
	Users                    []user.User
	Songs                    []enqueued_song.EnqueuedSong
	Devices                  []device.Device
	LoggedInUserRefreshToken internal_credentials.InternalCredentials
	ExternalCredentials      []external_credentials.ExternalCredentials
}{
	Rooms:                    rooms,
	Songs:                    songs,
	Users:                    users,
	Devices:                  devices,
	LoggedInUserRefreshToken: internalCredentials[0],
	ExternalCredentials:      externalCredentials,
}

type Seeder struct {
	Queryer i_queryer.IQueryer
}

func NewSeeder(queryer i_queryer.IQueryer) *Seeder {
	return &Seeder{
		Queryer: queryer,
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	for _, user1 := range users {
		if err := s.seedUser(ctx, &user1); err != nil {
			return err
		}
		devices1 := user1.Devices()
		for i := range devices1 {
			if err := s.seedDevice(ctx, &devices1[i], user1.Id()); err != nil {
				return err
			}
		}
	}

	for _, creds := range externalCredentials {
		if err := s.seedExternalCredentials(ctx, &creds); err != nil {
			return err
		}
	}

	for i := range songs {
		if err := s.seedEnqueuedSong(ctx, &songs[i]); err != nil {
			return err
		}
	}

	for i := range rooms {
		room1 := &rooms[i]
		if err := s.seedRoom(ctx, room1); err != nil {
			return err
		}
		songsInRoom := room1.EnqueuedSongs()
		for _, songInRoom := range songsInRoom {
			if err := s.seedEnqueuedSong(ctx, &songInRoom); err != nil {
				return err
			}
			songInRoomId := songInRoom.Id()
			userId := userIds[0]
			voteStatus := vote_status.Upvoted
			if err := s.seedUsersVotes(ctx, &userId, &songInRoomId, &voteStatus); err != nil {
				return err
			}
		}
		for _, memberId := range room1.Members() {
			var userRole *user_role.UserRole
			for _, u := range users {
				if u.Id() == memberId {
					userRole = u.Role()
					break
				}
			}
			if err := s.seedUserRoomData(ctx, room1.Id(), memberId, userRole); err != nil {
				return err
			}
		}
	}
	for _, usersRefreshToken := range internalCredentials {
		if err := s.seedInternalCredentials(ctx, usersRefreshToken); err != nil {
			return err
		}
	}

	return nil
}
