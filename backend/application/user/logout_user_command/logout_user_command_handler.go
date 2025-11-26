package logout_user_command

import (
	"context"
	"fmt"

	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
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
	err := c.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		if command.DeviceId == nil {
			retireAllTokensErr := c.refreshTokenRepository.RetireAllTokensByUserId(ctx, &userId, c.unitOfWork.GetQueryer())
			if retireAllTokensErr != nil {
				return custom_errors.NewCustomError("LogoutUserCommandHandler.RetireAllTokensByUserId",
					fmt.Sprintf("Couldn't retire users tokens for user id %s.", *userId.String()),
					retireAllTokensErr,
					custom_errors.Unexpected)
			}
			return nil
		}
		deviceId := command.DeviceId
		retireTokenWithDeviceId := c.refreshTokenRepository.RetireTokenByUserIdAndDeviceId(ctx, &userId, deviceId, c.unitOfWork.GetQueryer())
		if retireTokenWithDeviceId != nil {
			return custom_errors.NewCustomError("LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId",
				fmt.Sprintf("Couldn't retire users tokens for user id %s and device id %s.", *userId.String(), *deviceId.String()),
				retireTokenWithDeviceId,
				custom_errors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
