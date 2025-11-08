package daos

import (
	"time"

	"github.com/google/uuid"
)

type RoomDbDao struct {
	Id                   uuid.UUID `db:"id"`
	Name                 string    `db:"name"`
	Password             []byte    `db:"password"`
	QrCodeHash           []byte    `db:"qr_code_hash"`
	BoostCooldownSeconds *uint8    `db:"boost_cooldown_seconds"`
	CreatedAtUtc         time.Time `db:"created_at_utc"`
	LifespanSeconds      uint32    `db:"lifespan_seconds"`
}
