package persistance_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/get_room/daos"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) JoinRoomById(ctx context.Context, userId user.Id, roomId shared.RoomId, queryer contracts.IQueryer) error {
	args := m.Called(ctx, userId, roomId, queryer)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomIdByNameAndPassword(ctx context.Context, roomName, roomPassword string, queryer contracts.IQueryer) (*shared.RoomId, error) {
	args := m.Called(ctx, roomName, roomPassword, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shared.RoomId), args.Error(1)
}

func (m *MockRoomRepository) CreateRoom(ctx context.Context, roomParam *room.Room, queryer contracts.IQueryer) error {
	args := m.Called(ctx, roomParam, queryer)
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomByUserId(ctx context.Context, userId user.Id, queryer contracts.IQueryer) (*daos.GetRoomDao, error) {
	args := m.Called(ctx, userId, queryer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*daos.GetRoomDao), args.Error(1)
}

func (m *MockRoomRepository) CheckUserMembership(ctx context.Context, userId user.Id, queryer contracts.IQueryer) bool {
	args := m.Called(ctx, userId, queryer)
	return args.Bool(0)
}

func (m *MockRoomRepository) LeaveRoom(ctx context.Context, id user.Id, queryer contracts.IQueryer) error {
	args := m.Called(ctx, id, queryer)
	return args.Error(0)
}
