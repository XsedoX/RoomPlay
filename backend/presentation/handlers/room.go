package handlers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application"
	"xsedox.com/main/application/response"
	"xsedox.com/main/application/room"
	"xsedox.com/main/presentation/presentationErrors"
)

type RoomHandler struct {
	createRoomCommandHandler application.ICommandHandler[*room.CreateCommand]
}

func NewRoomHandler(createCommandHandler application.ICommandHandler[*room.CreateCommand]) *RoomHandler {
	return &RoomHandler{
		createRoomCommandHandler: createCommandHandler,
	}
}

// CreateRoom godoc
// @Summary      Create a new room
// @Description  Creates a new room in the system
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      room.CreateCommand	true  "Create Room"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.Error
// @Failure      401   {object}  response.Error
// @Failure      500   {object}  response.Error
// @Router       /api/v1/room [post]
// @Security BearerAuth
func (rh *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var cmd room.CreateCommand
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		presentationErrors.WriteJsonFailure(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := rh.createRoomCommandHandler.Handle(r.Context(), &cmd); err != nil {
		presentationErrors.WriteJsonFailure(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response.Ok(nil)); err != nil {
		presentationErrors.WriteJsonFailure(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
