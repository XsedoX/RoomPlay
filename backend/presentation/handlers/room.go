package handlers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/response"
	"xsedox.com/main/application/room"
	"xsedox.com/main/presentation/errors"
)

type RoomHandler struct {
	createCommandHandler *room.CreateCommandHandler
}

func NewRoomHandler(createCommandHandler *room.CreateCommandHandler) *RoomHandler {
	return &RoomHandler{
		createCommandHandler: createCommandHandler,
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
		w.WriteHeader(http.StatusBadRequest)
		if inerr := json.NewEncoder(w).Encode(response.Failure(err.Error())); inerr != nil {
			http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
		}
		return
	}

	if err := rh.createCommandHandler.Handle(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if inerr := json.NewEncoder(w).Encode(response.Failure(err.Error())); inerr != nil {
			http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	if inerr := json.NewEncoder(w).Encode(response.Ok(nil)); inerr != nil {
		http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
	}
	return
}
