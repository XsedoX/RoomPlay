package controllers

import (
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/user/data"
	"xsedox.com/main/presentation/response"
)

type UserController struct {
	getUserDataQueryHandler contracts.IQueryHandler[*data.UserQueryResponse]
}

func NewUserController(getUserDataQueryHandler contracts.IQueryHandler[*data.UserQueryResponse]) *UserController {
	return &UserController{getUserDataQueryHandler: getUserDataQueryHandler}
}

// GetUserData handles the HTTP request to retrieve user data.
// @Summary Retrieve user data
// @Description Get user data based on the provided context
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /user/data [get]
// @Security BearerAuth
func (userController *UserController) GetUserData(w http.ResponseWriter, r *http.Request) {
	userData, err := userController.getUserDataQueryHandler.Handle(r.Context())
	if err != nil {
		response.WriteJsonFailure(w, err.Error(), http.StatusBadRequest)
	}
	response.WriteJsonSuccess(w, userData, http.StatusOK)
}
