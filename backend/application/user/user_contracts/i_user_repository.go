package user_contracts

import (
	"context"

	application_contracts "github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IUserRepository interface {
	Add(ctx context.Context, user *user.User, queryer application_contracts.IQueryer) error
	CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer application_contracts.IQueryer) bool
	Update(ctx context.Context, user *user.User, queryer application_contracts.IQueryer) error
	GetUserByExternalId(ctx context.Context, externalId string, queryer application_contracts.IQueryer) (*user.User, error)
	GetUserById(ctx context.Context, id user.Id, queryer application_contracts.IQueryer) (*user.User, error)
}
