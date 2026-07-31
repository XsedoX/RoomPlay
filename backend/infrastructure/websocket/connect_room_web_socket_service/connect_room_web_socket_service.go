package connect_room_web_socket_service

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/hub"
)

type ConnectRoomWebSocketServiceRequest{

}

type IConnectRoomWebSocketService interface {
	ConnectWebSocket() error
}

type ConnectRoomWebSocketService struct {
	getRoomQueryHandler i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse]
	mainHub             hub.IHub
}

func NewConnectWebSocketService(
	getRoomQueryHandler i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse],
	mainHub hub.IHub,
) *ConnectRoomWebSocketService {
	return &ConnectRoomWebSocketService{
		getRoomQueryHandler: getRoomQueryHandler,
		mainHub:             mainHub,
	}
}
