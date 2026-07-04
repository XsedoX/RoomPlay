package i_user_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) bool
	Update(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error
	GetUserByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) (*user.User, error)
	GetUserById(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) (*user.User, error)
}
