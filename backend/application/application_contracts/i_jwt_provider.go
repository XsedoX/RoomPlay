package application_contracts

import (
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IJwtProvider interface {
	GenerateToken(userId user.Id) (string, error)
	ValidateTokenAndGetUserId(tokenString string) (*user.Id, error)
}
