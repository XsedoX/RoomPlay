package contracts

import (
	"context"

	"xsedox.com/main/domain/user"
)

type IUserRepository interface {
	AddUser(ctx context.Context, user *user.User, queryer IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer IQueryer) bool
	UpdateUser(ctx context.Context, user *user.User, queryer IQueryer) error
}
