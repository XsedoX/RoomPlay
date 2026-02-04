package i_unit_of_work

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
)

type IUnitOfWork interface {
	ExecuteTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	ExecuteRead(ctx context.Context, fn func(ctx context.Context) error) error
	GetQueryer() i_queryer.IQueryer
}
