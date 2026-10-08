package test_helpers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
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
	Data       json.RawMessage                          `json:"data"`
}

func ReadMessageForAction(
	t *testing.T,
	conn *websocket.Conn,
	expectedAction websocket_action.WebSocketOutgoingAction,
) json.RawMessage {
	t.Helper()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))
	defer func() {
		_ = conn.SetReadDeadline(time.Time{})
	}()

	for {
		_, message, err := conn.ReadMessage()
		require.NoError(t, err)

		var envelope struct {
			ActionName string `json:"actionName"`
		}
		require.NoError(t, json.Unmarshal(message, &envelope))

		if envelope.ActionName == string(expectedAction) {
			return message
		}
	}
}

type WebSocketPatchTestResponseWrapper struct {
	ActionName websocket_action.WebSocketOutgoingAction `json:"actionName"`
	Data       []struct {
		Op    string          `json:"op"`
		Path  string          `json:"path"`
		Value json.RawMessage `json:"value"`
	} `json:"data"`
}
