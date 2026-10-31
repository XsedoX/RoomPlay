package i_user_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IUserRepository interface {
	Create(ctx context.Context, user *user.User) (*user_id.UserId, error)
	CheckIfExistByExternalId(ctx context.Context, externalId string) (bool, error)
	Update(ctx context.Context, user *user.User) error
	GetByExternalId(ctx context.Context, externalId string) (*user.User, error)
	GetById(ctx context.Context, id user_id.UserId) (*user.User, error)
}
