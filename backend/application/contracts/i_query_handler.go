package contracts

import (
	"context"
)

type IQueryHandler[TResponse any] interface {
	Handle(ctx context.Context) (TResponse, error)
}
