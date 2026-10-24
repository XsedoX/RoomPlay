package user_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_room_data_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_vote_dao"
	"github.com/google/uuid"
)

type UserDao struct {
	Id         uuid.UUID                           `gorm:"primaryKey;type:uuid;" db:"id"`
	Name       string                              `gorm:"column:name" db:"name"`
	Surname    string                              `gorm:"column:surname" db:"surname"`
	RoomData   *user_room_data_dao.UserRoomDataDao `gorm:"foreignKey:UserId;references:Id"`
	Devices    []device_dao.DeviceDao              `gorm:"foreignKey:UserId;references:Id"`
	VotedSongs []user_vote_dao.UserVoteDao         `gorm:"foreignKey:UserId;references:Id"`
}

func (UserDao) TableName() string {
	return "users"
}

// ToDomain expects Devices and RoomData to have been loaded.
func (d UserDao) ToDomain() *user.User {
	devices := make([]device.Device, len(d.Devices))
	for i := range d.Devices {
		devices[i] = *d.Devices[i].ToDomain()
	}

	var roomId *room_id.RoomId
	var role *user_role.UserRole
	var boostUsedAtUtc *time.Time
	if d.RoomData != nil {
		membershipRoomId, membershipRole, boost := d.RoomData.ToDomain()
		roomId = new(membershipRoomId)
		role = new(membershipRole)
		boostUsedAtUtc = boost
	}

	return user.HydrateUser(
		user_id.UserId(d.Id),
		d.Name,
		d.Surname,
		role,
		roomId,
		devices,
		boostUsedAtUtc,
	)
}

// FromDomain maps the User aggregate, preserving loaded votes for the same user.
// Individual votes are not represented by the User domain object.
func (d *UserDao) FromDomain(value *user.User) {
	if d.Id != value.Id().ToUuid() {
		d.VotedSongs = nil
	}
	d.Id = value.Id().ToUuid()
	d.Name = value.FullName().Name()
	d.Surname = value.FullName().Surname()

	d.RoomData = nil
	if value.RoomId() != nil {
		d.RoomData = new(user_room_data_dao.UserRoomDataDao)
		d.RoomData.FromDomain(value.Id(), *value.RoomId(), *value.Role(), value.BoostUsedAtUtc())
	}

	devices := value.Devices()
	d.Devices = make([]device_dao.DeviceDao, len(devices))
	for i := range devices {
		d.Devices[i].FromDomain(&devices[i], value.Id())
	}
}
