package controllers

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/create_room_command"
	"xsedox.com/main/application/room/get_room_query"
	"xsedox.com/main/application/room/leave_room_command"
	"xsedox.com/main/infrastructure/validation"
	"xsedox.com/main/presentation/response"
)

type RoomController struct {
	createRoomCommandHandler          contracts.ICommandHandler[*create_room_command.CreateRoomCommand]
	getRoomQueryHandler               contracts.IQueryHandler[*get_room_query.GetRoomQueryResponse]
	getUserRoomMembershipQueryHandler contracts.IQueryHandler[*bool]
	leaveRoomCommandHandler           contracts.ICommandHandler[*leave_room_command.LeaveRoomCommand]
}

func NewRoomController(createRoomCommandHandler contracts.ICommandHandler[*create_room_command.CreateRoomCommand],
	getRoomQueryHandler contracts.IQueryHandler[*get_room_query.GetRoomQueryResponse],
	getUserRoomMembershipQueryHandler contracts.IQueryHandler[*bool],
	leaveRoomCommandHandler contracts.ICommandHandler[*leave_room_command.LeaveRoomCommand],
) *RoomController {
	return &RoomController{
		createRoomCommandHandler:          createRoomCommandHandler,
		getRoomQueryHandler:               getRoomQueryHandler,
		getUserRoomMembershipQueryHandler: getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler:           leaveRoomCommandHandler,
	}
}

// CreateRoom godoc
// @Summary      Join a new room
// @Description  Creates a new room in the system
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      create_room_command.CreateRoomCommand	true  "Join CreateRoomCommand"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router       /api/v1/room [post]
// @Security BearerAuth
func (rh *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var command create_room_command.CreateRoomCommand
	bodyDecodeErr := json.NewDecoder(r.Body).Decode(&command)
	if bodyDecodeErr != nil {
		response.WriteJsonFailure(w,
			"CreateRoomController.Decoding",
			"Problem with decoding request body",
			bodyDecodeErr.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}
	validationErr := validation.ValidatorInstance.Struct(command)
	if validationErr != nil {
		response.WriteJsonValidationFailure(w,
			"CreateRoom.Validation",
			r.URL.RequestURI(),
			validationErr)
		return
	}
	createRoomHandlerErr := rh.createRoomCommandHandler.Handle(r.Context(), &command)
	if createRoomHandlerErr != nil {
		response.WriteJsonApplicationFailure(w,
			createRoomHandlerErr,
			r.URL.RequestURI())
		return
	}
	response.WriteJsonSuccess(w, http.StatusCreated)
}

// GetRoom godoc
// @Summary      Get room data
// @Description  Gets details of a room.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router       /room [get]
// @Security BearerAuth
func (rh *RoomController) GetRoom(w http.ResponseWriter, r *http.Request) {
	getRoomQueryResponse, getRoomErr := rh.getRoomQueryHandler.Handle(r.Context())
	if getRoomErr != nil {
		response.WriteJsonApplicationFailure(w,
			getRoomErr,
			r.URL.RequestURI(),
		)
		return
	}
	response.WriteJsonSuccess(w, http.StatusOK, getRoomQueryResponse)
}

// CheckUserRoomMembership godoc
// @Summary      Checks if user is in any room.
// @Description  Returns true if user is in a room.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router       /room/membership [get]
// @Security BearerAuth
func (rh *RoomController) CheckUserRoomMembership(w http.ResponseWriter, r *http.Request) {
	handlerResponse, err := rh.getUserRoomMembershipQueryHandler.Handle(r.Context())
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI(),
		)
		return
	}
	response.WriteJsonSuccess(w, http.StatusOK, handlerResponse)
}

// LeaveRoom handles the HTTP request to leave a room.
// @Summary Makes a user leave a room
// @Description Used to leave a room
// @Tags room
// @Accept json
// @Produce json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router /room [delete]
// @Security BearerAuth
func (rh *RoomController) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	err := rh.leaveRoomCommandHandler.Handle(r.Context(), &leave_room_command.LeaveRoomCommand{})
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI())
	}
	response.WriteJsonSuccess(w, http.StatusOK)
}
