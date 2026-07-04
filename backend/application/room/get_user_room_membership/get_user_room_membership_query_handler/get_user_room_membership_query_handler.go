package get_user_room_membership_query_handler

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
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

func (g GetUserRoomMembershipQueryHandler) Handle(ctx context.Context) (*bool, error) {
	var result *bool
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application_helpers.NewMissingUserIdInContextError
	}
	err := g.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		_, err := g.roomRepository.GetRoomByUserId(ctx, *userId, g.unitOfWork.GetQueryer())
		if err != nil {
			tempResult := false
			result = &tempResult
			return custom_error.NewCustomError("GetUserRoomMembershipQueryHandler.GetRoomByUserId",
				fmt.Sprintf("Something went wrong with getting room by userId: %s", *userId.String()),
				err,
				custom_error_type.Unexpected)
		}
		tempResult := true
		result = &tempResult
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
