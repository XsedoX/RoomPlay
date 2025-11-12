package application

import (
	"context"

	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/domain/user"
)

const missingUserIdInContextErrorMessage string = "User id not found in context"

var (
	NewMissingUserIdInContextError = custom_errors.NewCustomError("GetUserIdFromContext.MissingUserContext",
		missingUserIdInContextErrorMessage,
		nil,
		custom_errors.Unauthorized)
)

func GetUserIdFromContext(ctx context.Context) (userId *user.Id, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user.Id)
	return userId, ok
}
