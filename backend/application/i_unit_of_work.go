package application

import "context"

type IUnitOfWork interface {
	ExecuteTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	ExecuteRead(ctx context.Context, fn func(ctx context.Context) error) error
}
