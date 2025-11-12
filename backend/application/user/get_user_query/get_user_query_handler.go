package get_user_query

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts2 "xsedox.com/main/application/user/contracts"
)

type GetUserQueryHandler struct {
	unitOfWork     contracts.IUnitOfWork
	userRepository contracts2.IUserRepository
}

func NewGetUserQueryHandler(unitOfWork contracts.IUnitOfWork,
	userRepository contracts2.IUserRepository) *GetUserQueryHandler {
	return &GetUserQueryHandler{
		unitOfWork:     unitOfWork,
		userRepository: userRepository}
}

func (handler GetUserQueryHandler) Handle(ctx context.Context) (*GetUserQueryResponse, error) {
	var response GetUserQueryResponse
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		user, er := handler.userRepository.GetUserById(ctx, *userId, handler.unitOfWork.GetQueryer())
		if er != nil {
			return er
		}
		response.Name = user.FullName().Name()
		response.Surname = user.FullName().Surname()
		return nil
	})
	if err != nil {
		return nil, custom_errors.NewCustomError(
			"GetUserQueryHandler.ExecuteTransaction",
			"Problem with executing transaction.",
			err,
			custom_errors.Unexpected)
	}
	return &response, nil
}
