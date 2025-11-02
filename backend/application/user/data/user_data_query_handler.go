package data

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
)

type UserQueryHandler struct {
	unitOfWork     contracts.IUnitOfWork
	userRepository contracts.IUserRepository
}

func NewUserQueryHandler(unitOfWork contracts.IUnitOfWork,
	userRepository contracts.IUserRepository) *UserQueryHandler {
	return &UserQueryHandler{
		unitOfWork:     unitOfWork,
		userRepository: userRepository}
}

func (handler UserQueryHandler) Handle(ctx context.Context) (*UserQueryResponse, error) {
	var response UserQueryResponse
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return nil, applicationErrors.NewApplicationError("", nil, applicationErrors.Unauthorized)
	}
	err := handler.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		user, er := handler.userRepository.GetUserById(ctx, *userId, handler.unitOfWork.GetQueryer())
		if er != nil {
			return er
		}
		response.Name = user.FullName().Name()
		response.Surname = user.FullName().Surname()
		response.Role = user.Role().String()
		roomId := user.RoomId().ToUuid()
		response.RoomId = roomId
		return nil
	})
	if err != nil {
		return nil, applicationErrors.NewApplicationError("", err, applicationErrors.Unexpected)
	}
	return &response, nil
}
