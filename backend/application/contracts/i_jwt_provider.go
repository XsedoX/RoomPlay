package contracts

import (
	"xsedox.com/main/domain/user"
)

type IJwtProvider interface {
	GenerateToken(userId user.Id) (string, error)
	ValidateTokenAndGetUserId(tokenString string) (*user.Id, error)
}
