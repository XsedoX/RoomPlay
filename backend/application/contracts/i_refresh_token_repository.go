package contracts

import (
	"context"

	"xsedox.com/main/domain/credentials"
)

type IRefreshTokenRepository interface {
	AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer IQueryer) error
	GetTokenByValue(ctx context.Context, value string, queryer IQueryer) (*credentials.RefreshToken, error)
	UpdateToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer IQueryer) error
}
