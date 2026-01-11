package contracts

import (
	"context"

	contracts2 "github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User, queryer contracts2.IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer contracts2.IQueryer) bool
	Update(ctx context.Context, user *user.User, queryer contracts2.IQueryer) error
	GetUserByExternalId(ctx context.Context, externalId string, queryer contracts2.IQueryer) (*user.User, error)
	GetUserById(ctx context.Context, id user.Id, queryer contracts2.IQueryer) (*user.User, error)
}
