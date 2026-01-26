package application_contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type IExternalCredentialsRepository interface {
	GetAccessTokenByUserId(ctx context.Context, userId user.Id, queryer IQueryer) (string, error)
	Grant(ctx context.Context, credentials *credentials.ExternalCredentials, queryer IQueryer) error
}
