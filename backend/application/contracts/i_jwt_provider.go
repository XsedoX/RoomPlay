package contracts

import (
	"xsedox.com/main/domain/shared"
)

type IJwtProvider interface {
	GenerateToken(userId shared.UserId) (string, error)
	ValidateTokenAndGetUserId(tokenString string) (*shared.UserId, error)
}
