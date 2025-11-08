package application

import (
	"context"

	"xsedox.com/main/domain/user"
)

const MissingUserIdInContextErrorMessage string = "User id not found in context"

func GetUserIdFromContext(ctx context.Context) (userId *user.Id, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user.Id)
	return userId, ok
}
