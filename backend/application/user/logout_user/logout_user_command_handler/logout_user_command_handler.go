package logout_user_command_handler

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command"
)

type LogoutUserCommandHandler struct {
	internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
	unitOfWork                    i_unit_of_work.IUnitOfWork
}

func NewLogoutUserCommandHandler(internalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository,
	unitOfWork i_unit_of_work.IUnitOfWork,
) *LogoutUserCommandHandler {
	return &LogoutUserCommandHandler{
		internalCredentialsRepository: internalCredentialsRepository,
		unitOfWork:                    unitOfWork,
	}
}

func (c LogoutUserCommandHandler) Handle(ctx context.Context, command *logout_user_command.LogoutUserCommand) error {
	userId := command.UserId
	err := c.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		if command.DeviceId == nil {
			retireAllTokensErr := c.internalCredentialsRepository.RetireAllTokensByUserId(ctx, &userId, c.unitOfWork.GetQueryer())
			if retireAllTokensErr != nil {
				return custom_error.NewCustomError("LogoutUserCommandHandler.RetireAllTokensByUserId",
					fmt.Sprintf("Couldn't retire users tokens for user id %s.", *userId.String()),
					retireAllTokensErr,
					custom_error_type.Unexpected)
			}
			return nil
		}
		deviceId := command.DeviceId
		retireTokenWithDeviceId := c.internalCredentialsRepository.RetireTokenByUserIdAndDeviceId(ctx, &userId, deviceId, c.unitOfWork.GetQueryer())
		if retireTokenWithDeviceId != nil {
			return custom_error.NewCustomError("LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId",
				fmt.Sprintf("Couldn't retire users tokens for user id %s and device id %s.", *userId.String(), *deviceId.String()),
				retireTokenWithDeviceId,
				custom_error_type.Unexpected)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
