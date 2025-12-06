package get_user_room_membership

import (
	"context"
	"fmt"

	"xsedox.com/main/application"
	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/application/room/contracts"
)

type GetUserRoomMembershipQueryHandler struct {
	roomRepository contracts.IRoomRepository
	unitOfWork     contracts2.IUnitOfWork
}

func NewGetUserRoomMembershipQueryHandler(roomRepository contracts.IRoomRepository,
	unitOfWork contracts2.IUnitOfWork) *GetUserRoomMembershipQueryHandler {
	return &GetUserRoomMembershipQueryHandler{
		roomRepository: roomRepository,
		unitOfWork:     unitOfWork,
	}
}
func (g GetUserRoomMembershipQueryHandler) Handle(ctx context.Context) (*bool, error) {
	var result *bool
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application.NewMissingUserIdInContextError
	}
	err := g.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		_, err := g.roomRepository.GetRoomByUserId(ctx, *userId, g.unitOfWork.GetQueryer())
		if err != nil {
			tempResult := false
			result = &tempResult
			return custom_errors.NewCustomError("GetUserRoomMembershipQueryHandler.GetRoomByUserId",
				fmt.Sprintf("Something went wrong with getting room by userId: %s", *userId.String()),
				err,
				custom_errors.Unexpected)
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
