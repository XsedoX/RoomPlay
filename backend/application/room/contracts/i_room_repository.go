package contracts

import (
	"context"

	contracts2 "github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IRoomRepository interface {
	CreateRoom(ctx context.Context, roomParam *room.Room, queryer contracts2.IQueryer) error
	GetRoomByUserId(ctx context.Context, userId user.Id, queryer contracts2.IQueryer) (*daos.GetRoomDao, error)
	CheckUserMembership(ctx context.Context, userId user.Id, queryer contracts2.IQueryer) bool
	LeaveRoom(ctx context.Context, id user.Id, queryer contracts2.IQueryer) error
	JoinRoomById(ctx context.Context, userId user.Id, roomId shared.RoomId, queryer contracts2.IQueryer) error
	GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer contracts2.IQueryer) (*shared.RoomId, error)
}
