package connect_room_web_socket_service

import (
	"context"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_publisher"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/hub"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/XsedoX/RoomPlay/presentation/response"
)

type IConnectRoomWebSocketService interface {
	ConnectWebSocket(ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		userId user_id.UserId,
		roomId room_id.RoomId,
	) error
}

type ConnectRoomWebSocketService struct {
	getRoomQueryHandler    i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse]
	mainHub                hub.IHub
	clientMessagePublisher client_message_publisher.IClientMessagePublisher
}

func NewConnectWebSocketService(
	getRoomQueryHandler i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse],
	mainHub hub.IHub,
	clientMessagePublisher client_message_publisher.IClientMessagePublisher,
) *ConnectRoomWebSocketService {
	return &ConnectRoomWebSocketService{
		getRoomQueryHandler:    getRoomQueryHandler,
		mainHub:                mainHub,
		clientMessagePublisher: clientMessagePublisher,
	}
}

func (c *ConnectRoomWebSocketService) ConnectWebSocket(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userId user_id.UserId,
	roomId room_id.RoomId,
) error {
	roomData, err := c.getRoomQueryHandler.Handle(ctx)
	if err != nil {
		return err
	}

	upgrader := c.mainHub.NewWebSocketUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := hub.NewClient(
		conn,
		userId,
		roomId,
		c.clientMessagePublisher,
		c.mainHub,
	)
	c.mainHub.RegisterClientToRoom(&hub.ClientRoomRequest{
		RoomId: room_id.RoomId(roomId),
		Client: client,
	})

	response := response.WebSocketSuccess[*get_room_query_response.GetRoomQueryResponse]{
		ActionName: websocket_action.GetRoomDataActionName,
		Success: response.Success[*get_room_query_response.GetRoomQueryResponse]{
			Data: roomData,
		},
	}
	c.mainHub.BroadcastToClientByConnectionId(&hub.ConnectionIdBroadcastRequest{
		client.ConnectionId(),
		response,
	})

	return nil
}
