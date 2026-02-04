package application_helpers

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

const missingUserIdInContextErrorMessage string = "User id not found in context"

var NewMissingUserIdInContextError = custom_error.NewCustomError("GetUserIdFromContext.MissingUserContext",
	missingUserIdInContextErrorMessage,
	nil,
	custom_error_type.Unauthorized)

func GetUserIdFromContext(ctx context.Context) (userId *user_id.UserId, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user_id.UserId)
	return userId, ok
}
