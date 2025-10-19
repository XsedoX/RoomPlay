package user

import (
	"context"

	"xsedox.com/main/application"
	"xsedox.com/main/domain/entities"
)

type CreateCommandHandler struct {
	userRepository IRepository
	unitOfWork     application.IUnitOfWork
}

func NewCreateCommandHandler(userRepository IRepository, unitOfWork application.IUnitOfWork) *CreateCommandHandler {
	return &CreateCommandHandler{userRepository, unitOfWork}
}
func (handler *CreateCommandHandler) Handle(ctx context.Context, command *CreateCommand) error {
	return handler.unitOfWork.Execute(ctx, func(ctx context.Context) error { return handler.handle(ctx, command) })
}
func (handler *CreateCommandHandler) handle(ctx context.Context, command *CreateCommand) error {
	device := entities.NewDevice(command.Device.Fingerprint, command.Device.FriendlyName, command.Device.DeviceType)
	user := entities.FirstLogin(command.ExternalId, command.Name, command.Surname, *device)
	return handler.userRepository.Create(ctx, user)
}
