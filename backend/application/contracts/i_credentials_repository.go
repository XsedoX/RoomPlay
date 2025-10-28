package contracts

import (
	"context"

	"xsedox.com/main/domain/credentials"
)

type ICredentialsRepository interface {
	Grant(ctx context.Context, credentials *credentials.External, queryer IQueryer) error
}
