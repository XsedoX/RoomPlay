package i_unit_of_work

import (
	"context"
)

type IUnitOfWork interface {
	ExecuteTransaction(ctx context.Context, fn func(context.Context) error) error
}
