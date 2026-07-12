package user_controller

import (
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_response"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

const (
	UserBasePath = "/user"
)

type UserController struct {
	getUserDataQueryHandler i_query_handler.IQueryHandler[*get_user_query_response.GetUserDataQueryResponse]
}

func NewUserController(getUserDataQueryHandler i_query_handler.IQueryHandler[*get_user_query_response.GetUserDataQueryResponse]) *UserController {
	return &UserController{
		getUserDataQueryHandler: getUserDataQueryHandler,
	}
}

// GetUserData handles the HTTP request to retrieve user data.
// @Summary Retrieve user data
// @Description Get user data based on the provided context
// @Tags user
// @Accept json
// @Produce json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router /user [get]
// @Security BearerAuth
func (userController *UserController) GetUserData(w http.ResponseWriter, r *http.Request) {
	userData, err := userController.getUserDataQueryHandler.Handle(r.Context())
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI())
	}
	response.WriteJsonSuccess(w, userData)
}
