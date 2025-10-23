package handlers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application"
	"xsedox.com/main/application/response"
	"xsedox.com/main/application/user"
	"xsedox.com/main/presentation/presentationErrors"
)

type UserHandler struct {
	createUserCommandHandler application.ICommandHandler[*user.LoginCommand]
}

func NewUserHandler(createUserHandler application.ICommandHandler[*user.LoginCommand]) *UserHandler {
	return &UserHandler{createUserCommandHandler: createUserHandler}
}

// LoginUser godoc
// @Summary      Logins a new user with data from OAuth
// @Description  Creates a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        room  body      user.LoginCommand	true  "Login user"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.Error
// @Failure      401   {object}  response.Error
// @Failure      500   {object}  response.Error
// @Router       /api/v1/user [post]
// @Security BearerAuth
func (handler *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(response.Failure("Authorization header required")); err != nil {
			presentationErrors.WriteJsonFailure(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}
}
