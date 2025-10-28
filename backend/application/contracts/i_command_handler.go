package contracts

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
)

type ICommandHandler[TCommand any] interface {
	Handle(ctx context.Context, command TCommand) *applicationErrors.ApplicationError
}

type ICommandHandlerWithResponse[TCommand any, TResponse any] interface {
	Handle(ctx context.Context, command TCommand) (TResponse, *applicationErrors.ApplicationError)
}
