package i_room_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IRoomRepository interface {
	Update(ctx context.Context, roomParam *room.Room) error
	Create(ctx context.Context, roomParam *room.Room) (*room_id.RoomId, error)
	GetUserMembership(ctx context.Context, userId user_id.UserId) *room_id.RoomId
	Leave(ctx context.Context, id user_id.UserId) error
	JoinById(ctx context.Context, userId user_id.UserId, roomId room_id.RoomId) error
	GetIdByNameAndPassword(ctx context.Context, roomName, roomPassword string) (*room_id.RoomId, error)
	GetEnqueuedSongAddedByByRoomIdEnqueuedSongId(ctx context.Context, roomId room_id.RoomId, enqueuedSongId enqueued_song_id.EnqueuedSongId) (string, error)
	GetByUserId(ctx context.Context, userId user_id.UserId) (*room.Room, error)
	GetById(ctx context.Context, roomId room_id.RoomId) (*room.Room, error)
}
