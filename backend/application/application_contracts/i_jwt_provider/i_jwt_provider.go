package i_jwt_provider

import (
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IJwtProvider interface {
	GenerateToken(userId user_id.UserId) (string, error)
	ValidateTokenAndGetUserId(tokenString string) (*user_id.UserId, error)
}
