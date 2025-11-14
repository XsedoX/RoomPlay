package logout_user_command

import (
	"context"

	"xsedox.com/main/application/contracts"
)

type LogoutUserCommandHandler struct {
	refreshTokenRepository contracts.IRefreshTokenRepository
	unitOfWork             contracts.IUnitOfWork
}

func NewLogoutUserCommandHandler(refreshTokenRepository contracts.IRefreshTokenRepository,
	unitOfWork contracts.IUnitOfWork) *LogoutUserCommandHandler {
	return &LogoutUserCommandHandler{
		refreshTokenRepository: refreshTokenRepository,
		unitOfWork:             unitOfWork,
	}
}

func (c LogoutUserCommandHandler) Handle(ctx context.Context, command *LogoutUserCommand) error {
	userId := command.UserId
	if command.DeviceId == nil {
		return c.refreshTokenRepository.RetireTokenByUserId(ctx, &userId, c.unitOfWork.GetQueryer())
	}
	deviceId := command.DeviceId
	return c.refreshTokenRepository.RetireTokenByUserIdAndDeviceId(ctx, &userId, deviceId, c.unitOfWork.GetQueryer())
}
