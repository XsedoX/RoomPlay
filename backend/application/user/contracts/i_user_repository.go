package contracts

import (
	"context"

	contracts2 "xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/user"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User, queryer contracts2.IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer contracts2.IQueryer) bool
	Update(ctx context.Context, user *user.User, queryer contracts2.IQueryer) error
	GetUserByExternalId(ctx context.Context, externalId string, queryer contracts2.IQueryer) (*user.User, error)
	GetUserById(ctx context.Context, id user.Id, queryer contracts2.IQueryer) (*user.User, error)
	LeaveRoom(ctx context.Context, id user.Id, queryer contracts2.IQueryer) error
}
