package seeder

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
)

var SeedData = struct {
	Rooms               []room.Room
	Users               []user.User
	Songs               []enqueued_song.EnqueuedSong
	Devices             []device.Device
	InternalCredentials []internal_credentials.InternalCredentials
	ExternalCredentials []external_credentials.ExternalCredentials
	UsersVotes          []UsersVotesStruct
}{
	Rooms:               rooms,
	Songs:               enqueuedSongs,
	Users:               users,
	Devices:             devices,
	InternalCredentials: internalCredentials,
	ExternalCredentials: externalCredentials,
	UsersVotes:          usersVotes,
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
	// NOTE: Insert rooms
	for _, roomToInsert := range SeedData.Rooms {
		err := s.seedRoom(ctx, &roomToInsert)
		if err != nil {
			return err
		}
	}
	// NOTE: Insert users
	for _, userToInsert := range SeedData.Users {
		err := s.seedUser(ctx, &userToInsert)
		if err != nil {
			return err
		}
	}
	// NOTE: Insert songs and enqueued songs
	for _, roomToInsert := range SeedData.Rooms {
		roomId := roomToInsert.Id()
		for _, enqueuedSongToInsert := range roomToInsert.AllSongs() {
			if err := s.seedEnqueuedSong(ctx, &enqueuedSongToInsert, &roomId); err != nil {
				return err
			}
		}
	}
	// NOTE: Insert external credentials
	for _, externalCredentialsToInsert := range SeedData.ExternalCredentials {
		err := s.seedExternalCredentials(ctx, &externalCredentialsToInsert)
		if err != nil {
			return err
		}
	}
	// NOTE: Insert devices
	for _, deviceOwner := range SeedData.Users {
		for _, deviceToInsert := range deviceOwner.Devices() {
			err := s.seedDevice(ctx, &deviceToInsert, deviceOwner.Id())
			if err != nil {
				return err
			}
		}
	}
	// NOTE: Insert internal credentials
	for _, internalCredentialsToInsert := range SeedData.InternalCredentials {
		if err := s.seedInternalCredentials(ctx, &internalCredentialsToInsert); err != nil {
			return err
		}
	}
	// NOTE: Insert votes
	for _, usersVotesToInsert := range SeedData.UsersVotes {
		for _, enqueuedSongId := range usersVotesToInsert.EnqueuedSongsIds {
			if err := s.seedUsersVotes(
				ctx,
				usersVotesToInsert.UserId,
				enqueuedSongId,
				usersVotesToInsert.VoteStatus,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
