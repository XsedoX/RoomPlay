package persistance

import (
	"context"

	"github.com/google/uuid"
	"xsedox.com/main/domain/entities"
)

type RoomRepository struct {
	queryer IQueryer
}

func NewRoomRepository(q IQueryer) *RoomRepository {
	return &RoomRepository{
		queryer: q,
	}
}

func (rr *RoomRepository) Create(ctx context.Context, room *entities.Room) error {
	_, err := rr.queryer.ExecContext(ctx,
		"INSERT INTO rooms (id, name, password, qr_code, boost_cooldown_seconds, created_at_utc, lifespan_seconds) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		uuid.UUID(room.GetId()),
		room.GetRoomName(),
		room.GetRoomPassword(),
		room.GetQrCode(),
		room.GetBoostCooldownSeconds(),
		room.GetCreatedAtUtc(),
		room.GetLifespanSeconds())
	if err != nil {
		return err
	}
	return nil
}
