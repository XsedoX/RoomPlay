package i_command_handler

import (
	"context"
)

type ICommandHandler[TCommand any] interface {
	Handle(ctx context.Context, command TCommand) error
}

type ICommandHandlerWithResponse[TCommand any, TResponse any] interface {
	Handle(ctx context.Context, command TCommand) (TResponse, error)
}
