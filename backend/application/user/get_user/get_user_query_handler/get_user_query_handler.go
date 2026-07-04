package get_user_query_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_response"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
)

type GetUserDataQueryHandler struct {
	unitOfWork     i_unit_of_work.IUnitOfWork
	userRepository i_user_repository.IUserRepository
}

func NewGetUserQueryHandler(unitOfWork i_unit_of_work.IUnitOfWork,
	userRepository i_user_repository.IUserRepository,
) *GetUserDataQueryHandler {
	return &GetUserDataQueryHandler{
		unitOfWork:     unitOfWork,
		userRepository: userRepository,
	}
}

func (handler GetUserDataQueryHandler) Handle(ctx context.Context) (*get_user_query_response.GetUserDataQueryResponse, error) {
	var response get_user_query_response.GetUserDataQueryResponse
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application_helpers.NewMissingUserIdInContextError
	}
	err := handler.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		user, er := handler.userRepository.GetUserById(ctx, *userId, handler.unitOfWork.GetQueryer())
		if er != nil {
			return custom_error.NewCustomError("NewGetUserQueryHandler.GetUserById",
				"Problem with creating a room.",
				er,
				custom_error_type.Unexpected)
		}
		response.Name = user.FullName().Name()
		response.Surname = user.FullName().Surname()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
