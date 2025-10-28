package login

import (
	"context"

	"xsedox.com/main/application/applicationErrors"
	"xsedox.com/main/application/contracts"
)

type UserCommandHandler struct {
	unitOfWork contracts.IUnitOfWork
}

func NewLoginUserCommandHandler(unitOfWork contracts.IUnitOfWork) *UserCommandHandler {
	return &UserCommandHandler{
		unitOfWork: unitOfWork,
	}
}
func (handler *UserCommandHandler) Handle(ctx context.Context, command *UserCommand) (*UserCommandResponse, *applicationErrors.ApplicationError) {
	var userCommandResponse UserCommandResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		return nil
	})
	return &userCommandResponse, applicationErrors.NewApplicationError("problem with executing transaction", err, applicationErrors.ErrInfrastructure)
}
