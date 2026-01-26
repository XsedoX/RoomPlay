package logout_user

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
)

type LogoutUserCommandHandler struct {
	refreshTokenRepository application_contracts.IRefreshTokenRepository
	unitOfWork             application_contracts.IUnitOfWork
}

func NewLogoutUserCommandHandler(refreshTokenRepository application_contracts.IRefreshTokenRepository,
	unitOfWork application_contracts.IUnitOfWork,
) *LogoutUserCommandHandler {
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
				return customerrors.NewCustomError("LogoutUserCommandHandler.RetireAllTokensByUserId",
					fmt.Sprintf("Couldn't retire users tokens for user id %s.", *userId.String()),
					retireAllTokensErr,
					customerrors.Unexpected)
			}
			return nil
		}
		deviceId := command.DeviceId
		retireTokenWithDeviceId := c.refreshTokenRepository.RetireTokenByUserIdAndDeviceId(ctx, &userId, deviceId, c.unitOfWork.GetQueryer())
		if retireTokenWithDeviceId != nil {
			return customerrors.NewCustomError("LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId",
				fmt.Sprintf("Couldn't retire users tokens for user id %s and device id %s.", *userId.String(), *deviceId.String()),
				retireTokenWithDeviceId,
				customerrors.Unexpected)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
