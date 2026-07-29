package get_user_room_membership_query_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
)

type GetUserRoomMembershipQueryHandler struct {
	roomRepository i_room_repository.IRoomRepository
	unitOfWork     i_unit_of_work.IUnitOfWork
}

func NewGetUserRoomMembershipQueryHandler(roomRepository i_room_repository.IRoomRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
) *GetUserRoomMembershipQueryHandler {
	return &GetUserRoomMembershipQueryHandler{
		roomRepository: roomRepository,
		unitOfWork:     unitOfWork,
	}
}

func (g GetUserRoomMembershipQueryHandler) Handle(ctx context.Context) (*room_id.RoomId, error) {
	var result *room_id.RoomId
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application_helpers.NewMissingUserIdInContextError("GetUserRoomMembershipQueryHandler.Handle")
	}
	err := g.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		roomId := g.roomRepository.GetUserMembership(ctx, *userId, g.unitOfWork.GetQueryer(ctx))
		result = roomId
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
