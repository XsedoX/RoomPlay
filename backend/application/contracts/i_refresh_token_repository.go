package contracts

import (
	"context"

	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
)

type IRefreshTokenRepository interface {
	AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer IQueryer) error
	GetTokenByValue(ctx context.Context, value string, queryer IQueryer) (*credentials.RefreshToken, error)
	RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user.Id, deviceId *user.DeviceId, queryer IQueryer) error
	RetireTokenByUserId(ctx context.Context, userId *user.Id, queryer IQueryer) error
}
