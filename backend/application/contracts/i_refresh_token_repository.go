package contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IRefreshTokenRepository interface {
	AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer IQueryer) error
	GetTokenByValue(ctx context.Context, value string, queryer IQueryer) (*credentials.RefreshToken, error)
	RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user.Id, deviceId *user.DeviceId, queryer IQueryer) error
	RetireAllTokensByUserId(ctx context.Context, userId *user.Id, queryer IQueryer) error
}
