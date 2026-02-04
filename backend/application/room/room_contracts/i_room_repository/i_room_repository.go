package i_room_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IRoomRepository interface {
	CreateRoom(ctx context.Context, roomParam *room.Room, queryer i_queryer.IQueryer) error
	GetRoomByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*get_room_dao.GetRoomDao, error)
	CheckUserMembership(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) bool
	LeaveRoom(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) error
	JoinRoomById(ctx context.Context, userId user_id.UserId, roomId room_id.RoomId, queryer i_queryer.IQueryer) error
	GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer i_queryer.IQueryer) (*room_id.RoomId, error)
}
