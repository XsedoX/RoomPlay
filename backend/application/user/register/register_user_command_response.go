package register

import "xsedox.com/main/domain/shared"

type UserCommandResponse struct {
	RefreshToken string
	AccessToken  string
	DeviceId     shared.DeviceId
}
