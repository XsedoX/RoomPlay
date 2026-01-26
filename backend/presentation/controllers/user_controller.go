package controllers

import (
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/user/get_user"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

const (
	UserBasePath = "/user"
)

type UserController struct {
	getUserDataQueryHandler application_contracts.IQueryHandler[*get_user.GetUserQueryResponse]
}

func NewUserController(getUserDataQueryHandler application_contracts.IQueryHandler[*get_user.GetUserQueryResponse]) *UserController {
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
	response.WriteJsonSuccess(w, http.StatusOK, userData)
}
