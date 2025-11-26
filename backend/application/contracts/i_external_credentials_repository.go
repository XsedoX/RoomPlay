package contracts

import (
	"context"

	"xsedox.com/main/domain/credentials"
)

type IExternalCredentialsRepository interface {
	Grant(ctx context.Context, credentials *credentials.External, queryer IQueryer) error
}
