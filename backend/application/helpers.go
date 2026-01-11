package application

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/domain/user"
)

const missingUserIdInContextErrorMessage string = "User id not found in context"

var NewMissingUserIdInContextError = customerrors.NewCustomError("GetUserIdFromContext.MissingUserContext",
	missingUserIdInContextErrorMessage,
	nil,
	customerrors.Unauthorized)

func GetUserIdFromContext(ctx context.Context) (userId *user.Id, ok bool) {
	value := ctx.Value(user.IdClaimContextKeyName)
	if value == nil {
		return nil, false
	}
	userId, ok = value.(*user.Id)
	return userId, ok
}
