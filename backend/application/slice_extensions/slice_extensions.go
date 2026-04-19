package slice_extensions

import (
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

func GetDeviceById(devices []device.Device, deviceId device_id.DeviceId) (*device.Device, bool) {
	for _, device := range devices {
		deviceId1 := device.Id()
		deviceId2 := deviceId
		if device_id.IdsEqual(&deviceId1, &deviceId2) {
			return &device, true
		}
	}
	return nil, false
}

func GetUserById(users []user.User, userId user_id.UserId) (*user.User, bool) {
	for _, user := range users {
		userId1 := user.Id()
		userId2 := userId
		if user_id.IdsEqual(&userId2, &userId1) {
			return &user, true
		}
	}
	return nil, false
}

func GetRoomById(rooms []room.Room, roomId room_id.RoomId) (*room.Room, bool) {
	for _, room := range rooms {
		roomId1 := room.Id()
		roomId2 := roomId
		if room_id.IdsEqual(&roomId1, &roomId2) {
			return &room, true
		}
	}
	return nil, false
}

func GetEnqueuedSongById(enqueuedSongs []enqueued_song.EnqueuedSong, enqueuedSongId enqueued_song_id.EnqueuedSongId) (*enqueued_song.EnqueuedSong, bool) {
	for _, enqueuedSong := range enqueuedSongs {
		enqueuedSongId1 := enqueuedSong.Id()
		enqueuedSongId2 := enqueuedSongId
		if enqueued_song_id.IdsEqual(&enqueuedSongId1, &enqueuedSongId2) {
			return &enqueuedSong, true
		}
	}
	return nil, false
}
