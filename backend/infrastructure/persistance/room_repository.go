package persistance

import (
	"context"

	"github.com/jackc/pgx/v5"
	"xsedox.com/domain/entities"
)

type RoomRepository struct {
	db  *pgx.Conn
	ctx *context.Context
}

func NewRoomRepository(db *pgx.Conn, ctx *context.Context) *RoomRepository {
	return &RoomRepository{
		db:  db,
		ctx: ctx,
	}
}

func (rp *RoomRepository) Create(room *entities.Room) error {
	_, err := rp.db.Query(*rp.ctx,
		"INSERT INTO rooms (id, name, password, user_id) VALUES ($1, $2, $3, $4)",
		room.GetId(),
		room.GetRoomName(),
		room.GetRoomPassword(),
		room.GetUserIds())
	if err != nil {
		return err
	}
	return nil
}
