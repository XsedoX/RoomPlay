package i_external_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IExternalCredentialsRepository interface {
	GetByUserId(ctx context.Context, userId user_id.UserId) (*external_credentials.ExternalCredentials, error)
	Upsert(ctx context.Context, credentials *external_credentials.ExternalCredentials) error
}
