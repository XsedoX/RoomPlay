package logout_command

import (
	"context"

	"xsedox.com/main/application/contracts"
)

type CommandHandler struct {
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
}

func NewLogoutRefreshTokenCommandHandler(refreshTokenRepository contracts.IRefreshTokenRepository,
	unitOfWork contracts.IUnitOfWork) *CommandHandler {
	return &CommandHandler{
		refreshTokenRepository: refreshTokenRepository,
		unitOfWork:             unitOfWork,
	}
}

func (c CommandHandler) Handle(ctx context.Context, command *Command) error {
	userId := command.UserId
	if command.DeviceId == nil {
		return c.refreshTokenRepository.RetireTokenByUserId(ctx, &userId, c.unitOfWork.GetQueryer())
	}
	deviceId := command.DeviceId
	return c.refreshTokenRepository.RetireTokenByUserIdAndDeviceId(ctx, &userId, deviceId, c.unitOfWork.GetQueryer())
}
