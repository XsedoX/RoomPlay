package application

import "context"

type IUnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}
