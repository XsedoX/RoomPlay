package room_contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IRoomRepository interface {
	CreateRoom(ctx context.Context, roomParam *room.Room, queryer application_contracts.IQueryer) error
	GetRoomByUserId(ctx context.Context, userId user.Id, queryer application_contracts.IQueryer) (*daos.GetRoomDao, error)
	CheckUserMembership(ctx context.Context, userId user.Id, queryer application_contracts.IQueryer) bool
	LeaveRoom(ctx context.Context, id user.Id, queryer application_contracts.IQueryer) error
	JoinRoomById(ctx context.Context, userId user.Id, roomId shared.RoomId, queryer application_contracts.IQueryer) error
	GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer application_contracts.IQueryer) (*shared.RoomId, error)
}
