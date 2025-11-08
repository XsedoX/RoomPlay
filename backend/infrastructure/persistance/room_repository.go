package persistance

import (
	"context"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/user"
)

type RoomRepository struct {
	encrypter contracts.IEncrypter
}

func NewRoomRepository(encrypter contracts.IEncrypter) *RoomRepository {
	return &RoomRepository{
		encrypter: encrypter,
	}
}

func (rr *RoomRepository) Create(ctx context.Context, roomParam *room.Room, queryer contracts.IQueryer) error {
	roomId := roomParam.Id()
	userId := roomParam.Members()[0]
	hashedPassword, err := rr.encrypter.HashAndSalt(roomParam.Password())
	if err != nil {
		return err
	}
	_, addRoomErr := queryer.ExecContext(ctx,
		`INSERT INTO rooms VALUES ($1::uuid, $2, $3::bytea, $4::bytea, $5, $6, $7);`,
		roomId.ToUuid(),
		roomParam.Name(),
		hashedPassword,
		[]byte(roomParam.QrCode()),
		roomParam.BoostCooldownSeconds(),
		roomParam.CreatedAtUtc().UTC(),
		roomParam.LifespanSeconds(),
	)
	if addRoomErr != nil {
		return addRoomErr
	}

	_, usersUpdateErr := queryer.ExecContext(ctx,
		`UPDATE users SET room_id = $1::uuid WHERE id = $2::uuid;`,
		roomId.ToUuid(),
		userId.ToUuid(),
	)
	if usersUpdateErr != nil {
		return usersUpdateErr
	}

	host := user.Host
	_, usersRolesErr := queryer.ExecContext(ctx,
		"INSERT INTO users_roles VALUES ($1, $2, $3);",
		roomId.ToUuid(),
		userId.ToUuid(),
		host.String(),
	)
	if usersRolesErr != nil {
		return usersRolesErr
	}
	return nil
}
