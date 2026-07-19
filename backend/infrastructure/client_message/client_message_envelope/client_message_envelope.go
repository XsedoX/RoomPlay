package client_message_envelope

import "github.com/XsedoX/RoomPlay/domain/user/user_id"

type ClientMessageEnvelope struct {
	UserId     user_id.UserId `json:"userId"`
	ActionName string         `json:"actionName"`
	Payload    []byte         `json:"payload"`
}
