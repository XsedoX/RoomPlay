package banned_user_dao

import (
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/google/uuid"
)

type BannedUserDao struct {
	RoomId  uuid.UUID        `gorm:"primaryKey;foreignKey:RoomId;references:Id"`
	UserId  uuid.UUID        `gorm:"primaryKey;foreignKey:UserId;references:Id"`
	UserDao user_dao.UserDao `gorm:"foreignKey:UserId;references:Id"`
}

func (BannedUserDao) TableName() string {
	return "banned_users"
}

func (d *BannedUserDao) FromDomain(value user_id.UserId, roomId room_id.RoomId) {
	*d = BannedUserDao{
		RoomId: roomId.ToUuid(),
		UserId: value.ToUuid(),
	}
}
