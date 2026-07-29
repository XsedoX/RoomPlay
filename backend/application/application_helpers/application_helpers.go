package application_helpers

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

const missingUserIdInContextErrorMessage string = "User id not found in context"

type UserIdContextKey string

const IdClaimContextKeyName UserIdContextKey = "user_id"

var NewMissingUserIdInContextError = func(id string) error {
	return application_error.NewApplicationError(id+".GetUserIdFromContext.MissingUserContext",
		missingUserIdInContextErrorMessage,
		nil,
		application_error_type.Unauthorized)
}

func GetUserIdFromContext(ctx context.Context) (userId *user_id.UserId, ok bool) {
	value := ctx.Value(IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user_id.UserId)
	return userId, ok
}
