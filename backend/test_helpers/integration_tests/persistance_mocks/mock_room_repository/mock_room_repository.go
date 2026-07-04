package mock_room_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) JoinRoomById(ctx context.Context, userId user_id.UserId, roomId room_id.RoomId, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, userId, roomId, queryer)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer i_queryer.IQueryer) (*room_id.RoomId, error) {
	args := m.Called(ctx, roomName, roomPassword, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*room_id.RoomId), args.Error(1)
}

func (m *MockRoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, roomParam, queryer)
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*get_room_dao.GetRoomDao, error) {
	args := m.Called(ctx, userId, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*get_room_dao.GetRoomDao), args.Error(1)
}

func (m *MockRoomRepository) CheckUserMembership(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) bool {
	args := m.Called(ctx, userId, queryer)
	return args.Bool(0)
}

func (m *MockRoomRepository) LeaveRoom(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, id, queryer)
	return args.Error(0)
}
