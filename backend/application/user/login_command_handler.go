package user

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/domain/entities"
)

type LoginCommandHandler struct {
	userRepository IRepository
	unitOfWork     application.IUnitOfWork
}

func NewLoginCommandHandler(userRepository IRepository, unitOfWork application.IUnitOfWork) *LoginCommandHandler {
	return &LoginCommandHandler{userRepository, unitOfWork}
}
func (handler *LoginCommandHandler) Handle(ctx context.Context, command *LoginCommand) error {
	return handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		device := entities.NewDevice(command.Device.Fingerprint, command.Device.DeviceType)
		user := entities.FirstLogin(command.ExternalId, command.Name, command.Surname, *device)
		return handler.userRepository.Create(ctx, user)
	})
}
