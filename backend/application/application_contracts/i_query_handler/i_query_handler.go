package i_query_handler

import (
	"context"
)

type IQueryHandlerWithRequest[TRequest any, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type IQueryHandler[TResponse any] interface {
	Handle(ctx context.Context) (TResponse, error)
}
