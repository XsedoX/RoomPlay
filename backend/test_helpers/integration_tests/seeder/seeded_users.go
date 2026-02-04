package seeder

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
)

var (
	userIds = []user_id.UserId{
		user_id.NewUserId(),
		user_id.NewUserId(),
		user_id.NewUserId(),
		user_id.NewUserId(),
		user_id.NewUserId(),
	}
	userRoles = []user_role.UserRole{
		user_role.Host,
		user_role.Member,
		user_role.Member,
		user_role.Member,
		user_role.Member,
	}
	users = []user.User{
		*user.HydrateUser(userIds[0],
			"name1",
			"surname1",
			&userRoles[0],
			&roomIds[1],
			[]device.Device{devices[4]},
			nil),
		*user.HydrateUser(userIds[1],
			"name2",
			"surname2",
			&userRoles[1],
			&roomIds[2],
			[]device.Device{devices[1]},
			nil),
		*user.HydrateUser(userIds[2],
			"name3",
			"surname3",
			&userRoles[2],
			&roomIds[2],
			[]device.Device{devices[2]},
			nil),
		*user.HydrateUser(userIds[3],
			"name4",
			"surname4",
			&userRoles[3],
			&roomIds[2],
			[]device.Device{devices[3]},
			nil),
		*user.HydrateUser(userIds[4],
			"name5",
			"surname5",
			&userRoles[4],
			nil,
			[]device.Device{devices[0]},
			nil),
	}
)

func (s *Seeder) seedUser(ctx context.Context, user *user.User) error {
	_, err := s.Queryer.ExecContext(ctx, `
		insert into users (id, name, surname)
		values ($1, $2, $3)
	`, user.Id(), user.FullName().Name(), user.FullName().Surname())
	if err != nil {
		return fmt.Errorf("failed to seed user: %w", err)
	}
	return nil
}

func (s *Seeder) seedUsersVotes(ctx context.Context, userId *user_id.UserId, enqueuedSongId *enqueued_song_id.EnqueuedSongId, voteStatus *vote_status.VoteStatus) error {
	_, err := s.Queryer.ExecContext(ctx, `
INSERT INTO users_votes (user_id, enqueued_song_id, vote_status)		
		VALUES ($1, $2, $3);`, userId, enqueuedSongId, voteStatus.String())
	if err != nil {
		return fmt.Errorf("failed to seed enqueued song %w", err)
	}
	return nil
}
