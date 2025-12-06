package contracts

import (
	"context"

	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/get_room/daos"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type IRoomRepository interface {
	CreateRoom(ctx context.Context, roomParam *room.Room, queryer contracts2.IQueryer) error
	GetRoomByUserId(ctx context.Context, userId user.Id, queryer contracts2.IQueryer) (*daos.GetRoomDao, error)
	CheckUserMembership(ctx context.Context, userId user.Id, queryer contracts2.IQueryer) bool
	LeaveRoom(ctx context.Context, id user.Id, queryer contracts2.IQueryer) error
	JoinRoomById(ctx context.Context, userId user.Id, roomId shared.RoomId, queryer contracts2.IQueryer) error
	GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer contracts2.IQueryer) (*shared.RoomId, error)
}
