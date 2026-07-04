package seeder

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/go-faker/faker/v4"
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
			faker.Name(),
			faker.LastName(),
			&userRoles[0],
			&roomIds[1],
			[]device.Device{devices[4]},
			nil),
		*user.HydrateUser(userIds[1],
			faker.Name(),
			faker.LastName(),
			&userRoles[1],
			&roomIds[2],
			[]device.Device{devices[1]},
			nil),
		*user.HydrateUser(userIds[2],
			faker.Name(),
			faker.LastName(),
			&userRoles[2],
			&roomIds[2],
			[]device.Device{devices[2]},
			nil),
		*user.HydrateUser(userIds[3],
			faker.Name(),
			faker.LastName(),
			&userRoles[3],
			&roomIds[2],
			[]device.Device{devices[3]},
			nil),
		*user.HydrateUser(userIds[4],
			faker.Name(),
			faker.LastName(),
			nil,
			nil,
			[]device.Device{devices[0]},
			nil),
	}
)

func (s *Seeder) seedUser(ctx context.Context, user *user.User) error {
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users 
		(
			id,
			name,
			surname
		)
		VALUES (
			$1, $2, $3
		)
`, user.Id().ToUuid(),
		user.FullName().Name(),
		user.FullName().Surname(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed user: %w", err)
	}
	// NOTE: Not every user is in a room
	if user.RoomId() == nil {
		return nil
	}
	_, err = s.Queryer.ExecContext(ctx, `
		INSERT INTO users_room_data
		(
			user_id,
			room_id,
			boost_used_at_utc,
			role
		)
		VALUES
		(
			$1, $2, $3, $4
		)
		`,
		user.Id().ToUuid(),
		user.RoomId().ToUuid(),
		user.BoostUsedAtUtc(),
		user.Role().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed users_room_data: %w", err)
	}
	return nil
}
