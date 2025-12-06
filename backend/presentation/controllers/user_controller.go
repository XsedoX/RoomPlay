package controllers

import (
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/user/get_user"
	"xsedox.com/main/presentation/response"
)

type UserController struct {
	getUserDataQueryHandler contracts.IQueryHandler[*get_user.GetUserQueryResponse]
}

func NewUserController(getUserDataQueryHandler contracts.IQueryHandler[*get_user.GetUserQueryResponse]) *UserController {
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
