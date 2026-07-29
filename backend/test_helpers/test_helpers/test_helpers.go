package test_helpers

import (
	"context"
	"encoding/json"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
)

func AddUserIdToContext(ctx context.Context) (user_id.UserId, context.Context) {
	userId := user_id.NewUserId()
	ctx = context.WithValue(ctx, application_helpers.IdClaimContextKeyName, &userId)
	return userId, ctx
}

type TestResponseWrapper[T any] struct {
	Data T                         `json:"data"`
	Meta page_meta_dto.PageMetaDto `json:"meta"`
}

type WebSocketTestResponseWrapper struct {
	ActionName websocket_action.WebSocketOutgoingAction `json:"actionName"`
	Data       []struct {
		Op    string          `json:"op"`
		Path  string          `json:"path"`
		Value json.RawMessage `json:"value"`
	} `json:"data"`
}
