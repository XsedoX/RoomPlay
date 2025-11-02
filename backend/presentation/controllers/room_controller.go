package controllers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/join"
	"xsedox.com/main/presentation/response"
)

type RoomController struct {
	createRoomCommandHandler contracts.ICommandHandler[*join.RoomCommand]
}

func NewRoomController(createCommandHandler contracts.ICommandHandler[*join.RoomCommand]) *RoomController {
	return &RoomController{
		createRoomCommandHandler: createCommandHandler,
	}
}

// CreateRoom godoc
// @Summary      Join a new room
// @Description  Creates a new room in the system
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      join.RoomCommand	true  "Join Room"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.Error
// @Failure      401   {object}  response.Error
// @Failure      500   {object}  response.Error
// @Router       /api/v1/room [post]
// @Security BearerAuth
func (rh *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var cmd join.RoomCommand
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		response.WriteJsonFailure(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := rh.createRoomCommandHandler.Handle(r.Context(), &cmd); err != nil {
		response.WriteJsonFailure(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response.Ok(nil)); err != nil {
		response.WriteJsonFailure(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
