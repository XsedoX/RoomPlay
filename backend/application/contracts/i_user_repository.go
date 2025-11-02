package contracts

import (
	"context"

	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User, queryer IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer IQueryer) bool
	Update(ctx context.Context, user *user.User, queryer IQueryer) error
	GetUserByExternalId(ctx context.Context, externalId string, queryer IQueryer) (*user.User, error)
	GetUserById(ctx context.Context, id shared.UserId, queryer IQueryer) (*user.User, error)
}
