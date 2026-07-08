package application_helpers

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

const missingUserIdInContextErrorMessage string = "User id not found in context"

var NewMissingUserIdInContextError = application_error.NewApplicationError("GetUserIdFromContext.MissingUserContext",
	missingUserIdInContextErrorMessage,
	nil,
	application_error_type.Unauthorized)

func GetUserIdFromContext(ctx context.Context) (userId *user_id.UserId, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user_id.UserId)
	return userId, ok
}
