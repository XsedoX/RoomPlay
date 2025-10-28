package application

import (
	"context"

	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

func GetUserIdFromContext(ctx context.Context) (userId *shared.UserId, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*shared.UserId)
	return userId, ok
}
