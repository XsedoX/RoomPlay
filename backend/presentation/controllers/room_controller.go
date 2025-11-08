package controllers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/create_command"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/infrastructure/validation"
	"xsedox.com/main/presentation/response"
)

type RoomController struct {
	createRoomCommandHandler contracts.ICommandHandlerWithResponse[*create_command.CreateRoomCommand, *shared.RoomId]
}

func NewRoomController(createRoomCommandHandler contracts.ICommandHandlerWithResponse[*create_command.CreateRoomCommand, *shared.RoomId]) *RoomController {
	return &RoomController{
		createRoomCommandHandler: createRoomCommandHandler,
	}
}

// CreateRoom godoc
// @Summary      Join a new room
// @Description  Creates a new room in the system
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      create_command.CreateRoomCommand	true  "Join CreateRoomCommand"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router       /api/v1/room [post]
// @Security BearerAuth
func (rh *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var command create_command.CreateRoomCommand
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		response.WriteJsonFailure(w,
			"CreateRoomController.Decoding",
			"Problem with decoding request body",
			err.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}
	err = validation.ValidatorInstance.Struct(command)
	if err != nil {
		response.WriteJsonValidationFailure(w,
			"CreateRoom.Validation",
			r.URL.RequestURI(),
			err)
		return
	}
	resp, err := rh.createRoomCommandHandler.Handle(r.Context(), &command)
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI())
		return
	}
	response.WriteJsonSuccess(w, resp, http.StatusCreated)
}
