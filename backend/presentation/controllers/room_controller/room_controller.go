package room_controller

import (
	"encoding/json"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_publisher"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/i_hub"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/main_hub"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/room_hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket_requests/client_room_request"
	"github.com/XsedoX/RoomPlay/presentation/response"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

const (
	RoomBasePath           = "/room"
	RoomMembershipBasePath = "/membership"
	JoinRoomPasswordPath   = "/join/password"
)

type RoomController struct {
	createRoomCommandHandler          i_command_handler.ICommandHandlerWithResponse[*create_room_command.CreateRoomCommand, *room_id.RoomId]
	getRoomQueryHandler               i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse]
	getUserRoomMembershipQueryHandler i_query_handler.IQueryHandler[*bool]
	leaveRoomCommandHandler           i_command_handler.ICommandHandler[*leave_room_command.LeaveRoomCommand]
	joinRoomCommandHandler            i_command_handler.ICommandHandler[*join_room_password_command.JoinRoomPasswordCommand]
	mainHub                           i_hub.IHub
	clientMessagePublisher            client_message_publisher.IClientMessagePublisher
}

func NewRoomController(createRoomCommandHandler i_command_handler.ICommandHandlerWithResponse[*create_room_command.CreateRoomCommand, *room_id.RoomId],
	getRoomQueryHandler i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse],
	getUserRoomMembershipQueryHandler i_query_handler.IQueryHandler[*bool],
	leaveRoomCommandHandler i_command_handler.ICommandHandler[*leave_room_command.LeaveRoomCommand],
	joinRoomCommandHandler i_command_handler.ICommandHandler[*join_room_password_command.JoinRoomPasswordCommand],
	mainHub i_hub.IHub,
	clientMessagePublisher client_message_publisher.IClientMessagePublisher,
) *RoomController {
	return &RoomController{
		createRoomCommandHandler:          createRoomCommandHandler,
		clientMessagePublisher:            clientMessagePublisher,
		getRoomQueryHandler:               getRoomQueryHandler,
		getUserRoomMembershipQueryHandler: getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler:           leaveRoomCommandHandler,
		joinRoomCommandHandler:            joinRoomCommandHandler,
		mainHub:                           mainHub,
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
	var command create_room_command.CreateRoomCommand
	bodyDecodeErr := json.NewDecoder(r.Body).Decode(&command)
	if bodyDecodeErr != nil {
		response.WriteJsonDecodingFailure(
			w,
			"CreateRoomController.Decoding",
			bodyDecodeErr,
			r.URL.RequestURI(),
		)
		return
	}
	validationErr := setup_validation.ValidatorInstance.Struct(command)
	if validationErr != nil {
		response.WriteJsonValidationFailure(w,
			"CreateRoom.Validation",
			r.URL.RequestURI(),
			validationErr)
		return
	}
	roomId, createRoomHandlerErr := rh.createRoomCommandHandler.Handle(r.Context(), &command)
	if createRoomHandlerErr != nil {
		response.WriteJsonApplicationFailure(w,
			createRoomHandlerErr,
			r.URL.RequestURI())
		return
	}
	response.WriteJsonCreated(w, roomId.ToUuid())
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
	response.WriteJsonSuccess(w, getRoomQueryResponse)
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
	response.WriteJsonSuccess(w, handlerResponse)
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
	var command join_room_password_command.JoinRoomPasswordCommand
	bodyDecodeErr := json.NewDecoder(r.Body).Decode(&command)
	if bodyDecodeErr != nil {
		response.WriteJsonDecodingFailure(
			w,
			"JoinRoomController.Password.Decoding",
			bodyDecodeErr,
			r.URL.RequestURI(),
		)
		return
	}
	validationErr := setup_validation.ValidatorInstance.Struct(command)
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

func (rh *RoomController) UpgradeToWebSockets(w http.ResponseWriter, r *http.Request) {
	userId, ok := application_helpers.GetUserIdFromContext(r.Context())
	if !ok {
		response.WriteJsonApplicationFailure(w, application_helpers.NewMissingUserIdInContextError, r.URL.RequestURI())
	}

	type roomIdParam struct {
		RoomId uuid.UUID `schema:"roomId" validate:"required,uuid"`
	}
	var roomIdParamValue roomIdParam
	decoder := schema.NewDecoder()
	paramsDecodeErr := decoder.Decode(&roomIdParamValue, r.URL.Query())
	if paramsDecodeErr != nil {
		response.WriteJsonDecodingFailure(
			w,
			"UpgradeToWebSockets.Decoding",
			paramsDecodeErr,
			r.URL.RequestURI(),
		)
		return
	}

	upgrader := main_hub.NewWebSocketUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.WriteJsonApplicationFailure(w,
			err,
			r.URL.RequestURI())
		return
	}
	client := room_hub.NewClient(
		conn,
		*userId,
		rh.clientMessagePublisher,
	)
	go client.WritePump()
	go client.ReadPump()
	rh.mainHub.RegisterClientToRoom(&client_room_request.ClientRoomRequest{
		RoomId: room_id.RoomId(roomIdParamValue.RoomId),
		Client: client,
	})
}
