package i_user_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string) (bool, error)
	Update(ctx context.Context, user *user.User) error
	GetUserByExternalId(ctx context.Context, externalId string) (*user.User, error)
	GetUserById(ctx context.Context, id user_id.UserId) (*user.User, error)
}
