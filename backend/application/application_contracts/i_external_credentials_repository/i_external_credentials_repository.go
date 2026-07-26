package i_external_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IExternalCredentialsRepository interface {
	GetAccessTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error)
	Grant(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error
	GetRefreshTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error)
	GetExternalCredentialsByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*external_credentials.ExternalCredentials, error)
	UpdateExternalCredentials(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error
}
