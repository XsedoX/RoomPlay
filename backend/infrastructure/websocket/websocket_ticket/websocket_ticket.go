package web_socket_ticket

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type WebSocketTicket struct {
	Token  string         `json:"token"`
	UserId user_id.UserId `json:"userId"`
	RoomId room_id.RoomId `json:"roomId"`
}

func NewWebSocketTicket(userId user_id.UserId, roomId room_id.RoomId, encrypter i_encrypter.IEncrypter) *WebSocketTicket {
	return &WebSocketTicket{
		Token:  encrypter.NewRandomString(),
		UserId: userId,
		RoomId: roomId,
	}
}
