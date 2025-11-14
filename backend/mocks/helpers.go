package mocks

import (
	"context"

	"github.com/google/uuid"
	"xsedox.com/main/domain/user"
)

func AddUserIdToContext(ctx context.Context) (user.Id, context.Context) {
	userId := user.Id(uuid.New())
	ctx = context.WithValue(ctx, user.IdClaimContextKeyName, &userId)
	return userId, ctx
}
