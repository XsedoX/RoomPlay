package persistance

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/domain/entities"
)

type RoomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (rr *RoomRepository) Create(ctx context.Context, room *entities.Room) error {
	queryer := GetQueryerFromContext(ctx, rr.db)

	_, err := queryer.QueryxContext(ctx,
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
