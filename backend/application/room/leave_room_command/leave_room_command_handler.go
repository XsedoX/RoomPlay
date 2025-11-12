package leave_room_command

import (
	"context"

	"xsedox.com/main/application"
	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/application/user/contracts"
)

type LeaveRoomCommandHandler struct {
	userRepository contracts.IUserRepository
	unitOfWork     contracts2.IUnitOfWork
}

func (l LeaveRoomCommandHandler) Handle(ctx context.Context, _ *LeaveRoomCommand) error {
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return application.NewMissingUserIdInContextError
	}
	err := l.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		err := l.userRepository.LeaveRoom(ctx, *userId, l.unitOfWork.GetQueryer())
		if err != nil {
			return custom_errors.NewCustomError("LeaveRoomCommandHandler.LeaveRoom",
				"Couldn't leave room.",
				err,
				custom_errors.Unexpected)
		}
		return nil
	})
	return err
}

func NewLeaveRoomCommandHandler(userRepository contracts.IUserRepository, unitOfWork contracts2.IUnitOfWork) *LeaveRoomCommandHandler {
	return &LeaveRoomCommandHandler{userRepository: userRepository, unitOfWork: unitOfWork}
}
