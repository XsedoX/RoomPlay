package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/room/create_room"
	"github.com/XsedoX/RoomPlay/application/room/get_room"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password"
	"github.com/XsedoX/RoomPlay/application/room/leave_room"
	"github.com/XsedoX/RoomPlay/infrastructure/validation"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

const (
	RoomBasePath           = "/room"
	RoomMembershipBasePath = "/membership"
	JoinRoomPasswordPath   = "/join/password"
)

type RoomController struct {
	createRoomCommandHandler          application_contracts.ICommandHandler[*create_room.CreateRoomCommand]
	getRoomQueryHandler               application_contracts.IQueryHandler[*get_room.GetRoomQueryResponse]
	getUserRoomMembershipQueryHandler application_contracts.IQueryHandler[*bool]
	leaveRoomCommandHandler           application_contracts.ICommandHandler[*leave_room.LeaveRoomCommand]
	joinRoomCommandHandler            application_contracts.ICommandHandler[*join_room_password.JoinRoomPasswordCommand]
}

func NewRoomController(createRoomCommandHandler application_contracts.ICommandHandler[*create_room.CreateRoomCommand],
	getRoomQueryHandler application_contracts.IQueryHandler[*get_room.GetRoomQueryResponse],
	getUserRoomMembershipQueryHandler application_contracts.IQueryHandler[*bool],
	leaveRoomCommandHandler application_contracts.ICommandHandler[*leave_room.LeaveRoomCommand],
	joinRoomCommandHandler application_contracts.ICommandHandler[*join_room_password.JoinRoomPasswordCommand],
) *RoomController {
	return &RoomController{
		createRoomCommandHandler:          createRoomCommandHandler,
		getRoomQueryHandler:               getRoomQueryHandler,
		getUserRoomMembershipQueryHandler: getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler:           leaveRoomCommandHandler,
		joinRoomCommandHandler:            joinRoomCommandHandler,
	}
}

// CreateRoom godoc
// @Summary      Join a new room
// @Description  Creates a new room in the system
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      create_room.CreateRoomCommand	true  "Join CreateRoomCommand"
// @Success      201   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router       /api/v1/room [post]
// @Security BearerAuth
func (rh *RoomController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var command create_room.CreateRoomCommand
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
	err := rh.leaveRoomCommandHandler.Handle(r.Context(), &leave_room.LeaveRoomCommand{})
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI())
	}
	response.WriteJsonSuccess(w, http.StatusOK)
}

// JoinRoomPassword handles the HTTP request to join a room.
// @Summary Makes a user join a room
// @Description Used to join a room
// @Tags room
// @Accept json
// @Produce json
// @Success      200   {object}  response.Success
// @Failure      400   {object}  response.ProblemDetails
// @Failure      401   {object}  response.ProblemDetails
// @Failure      500   {object}  response.ProblemDetails
// @Router /room/join [put]
// @Security BearerAuth
func (rh *RoomController) JoinRoomPassword(w http.ResponseWriter, r *http.Request) {
	var command join_room_password.JoinRoomPasswordCommand
	bodyDecodeErr := json.NewDecoder(r.Body).Decode(&command)
	if bodyDecodeErr != nil {
		response.WriteJsonFailure(w,
			"JoinRoomController.Password.Decoding",
			"Problem with decoding request body",
			bodyDecodeErr.Error(),
			r.URL.RequestURI(),
			http.StatusBadRequest)
		return
	}
	validationErr := validation.ValidatorInstance.Struct(command)
	if validationErr != nil {
		response.WriteJsonValidationFailure(w,
			"JoinRoomPassword.Validation",
			r.URL.RequestURI(),
			validationErr)
		return
	}
	joinRoomCommandHandlerErr := rh.joinRoomCommandHandler.Handle(r.Context(), &command)
	if joinRoomCommandHandlerErr != nil {
		response.WriteJsonApplicationFailure(w,
			joinRoomCommandHandlerErr,
			r.URL.RequestURI())
		return
	}
	response.WriteJsonNoContent(w)
}
