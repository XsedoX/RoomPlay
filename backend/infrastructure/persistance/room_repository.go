package persistance

import (
	"context"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/room"
)

type RoomRepository struct {
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (rr *RoomRepository) Join(ctx context.Context, room *room.Room, queryer contracts.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		"INSERT INTO rooms (id, name, password, qr_code, boost_cooldown_seconds, created_at_utc, lifespan_seconds) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		uuid.UUID(room.GetId()),
		room.GetRoomName(),
		room.GetRoomPassword(),
		room.GetQrCode(),
		room.GetBoostCooldownSeconds(),
		room.GetCreatedAtUtc(),
		room.GetLifespanSeconds())
	return err
}
