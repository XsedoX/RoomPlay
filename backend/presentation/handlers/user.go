package handlers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/response"
	"xsedox.com/main/application/user"
	"xsedox.com/main/presentation/errors"
)

type UserHandler struct {
	createUserCommandHandler *user.CreateCommandHandler
}

func NewUserHandler(createUserHandler *user.CreateCommandHandler) *UserHandler {
	return &UserHandler{createUserCommandHandler: createUserHandler}
}

// LoginUser godoc
// @Summary      Logins a new user with data from OAuth
// @Description  Creates a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        room  body      user.CreateCommand	true  "Login user"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.Error
// @Failure      401   {object}  response.Error
// @Failure      500   {object}  response.Error
// @Router       /api/v1/user [post]
// @Security BearerAuth
func (handler *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var cmd user.CreateCommand
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if inerr := json.NewEncoder(w).Encode(response.Failure(err.Error())); inerr != nil {
			http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	if err := handler.createUserCommandHandler.Handle(r.Context(), &cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if inerr := json.NewEncoder(w).Encode(response.Failure(err.Error())); inerr != nil {
			http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if inerr := json.NewEncoder(w).Encode(response.Ok(nil)); inerr != nil {
		http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
	}
	return
}
