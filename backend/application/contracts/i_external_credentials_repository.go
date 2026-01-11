package contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/credentials"
)

type IExternalCredentialsRepository interface {
	Grant(ctx context.Context, credentials *credentials.External, queryer IQueryer) error
}
