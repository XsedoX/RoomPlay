package user_room_data_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/google/uuid"
)

type UserRoomDataDao struct {
	UserId         uuid.UUID  `gorm:"column:user_id;primaryKey;foreignKey:UserId;references:Id" db:"user_id"`
	RoomId         uuid.UUID  `gorm:"column:room_id;primaryKey;foreignKey:UserId;references:Id" db:"room_id"`
	Role           string     `gorm:"column:role;type:user_role;not null;default:'member'" db:"role"`
	BoostUsedAtUtc *time.Time `gorm:"column:boost_used_at_utc" db:"boost_used_at_utc"`
}

func (UserRoomDataDao) TableName() string {
	return "users_room_data"
}

func (d *UserRoomDataDao) FromDomain(userId user_id.UserId, roomId room_id.RoomId, role user_role.UserRole, boostUsedAtUtc *time.Time) {
	*d = UserRoomDataDao{
		UserId: userId.ToUuid(),
		RoomId: roomId.ToUuid(),
		Role:   *role.String(),
	}
	if boostUsedAtUtc != nil {
		d.BoostUsedAtUtc = new(boostUsedAtUtc.UTC())
	}
}
