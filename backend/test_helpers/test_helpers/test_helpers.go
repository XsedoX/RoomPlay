package test_helpers

import (
	"context"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

func AddUserIdToContext(ctx context.Context) (user_id.UserId, context.Context) {
	userId := user_id.NewUserId()
	ctx = context.WithValue(ctx, user.IdClaimContextKeyName, &userId)
	return userId, ctx
}

func PtrString(s string) *string     { return &s }
func PtrTime(t time.Time) *time.Time { return &t }

type TestResponseWrapper[T any] struct {
	Data T `json:"data"`
}
