package logout_user_command_handler

import (
	"context"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
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
				return application_error.NewApplicationError("LogoutUserCommandHandler.RetireAllTokensByUserId",
					fmt.Sprintf("Couldn't retire users tokens for user id %s.", *userId.String()),
					retireAllTokensErr,
					application_error_type.Unexpected)
			}
			return nil
		}
		deviceId := command.DeviceId
		userSession := user_session.NewUserSession(userId, *deviceId)
		retireTokenWithDeviceId := c.internalCredentialsRepository.RetireTokenByUserSession(ctx, *userSession, c.unitOfWork.GetQueryer())
		if retireTokenWithDeviceId != nil {
			return application_error.NewApplicationError("LogoutUserCommandHandler.RetireTokenByUserIdAndDeviceId",
				fmt.Sprintf("Couldn't retire users tokens for user id %s and device id %s.", *userId.String(), *deviceId.String()),
				retireTokenWithDeviceId,
				application_error_type.Unexpected)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
