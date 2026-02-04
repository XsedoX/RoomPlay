package user_session

import (
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type UserSession struct {
	userId   user_id.UserId
	deviceId device_id.DeviceId
}

func NewUserSession(userId user_id.UserId, deviceId device_id.DeviceId) *UserSession {
	return &UserSession{
		deviceId: deviceId,
		userId:   userId,
	}
}

func (us *UserSession) UserId() user_id.UserId {
	return us.userId
}

func (us *UserSession) DeviceId() device_id.DeviceId {
	return us.deviceId
}
