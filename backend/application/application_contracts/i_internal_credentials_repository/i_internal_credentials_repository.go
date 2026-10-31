package i_internal_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IInternalCredentialsRepository interface {
	Update(ctx context.Context, refreshToken *internal_credentials.InternalCredentials) error
	GetByValue(ctx context.Context, value string) (*internal_credentials.InternalCredentials, error)
	RetireByUserSession(ctx context.Context, userSession user_session.UserSession) error
	RetireAllByUserId(ctx context.Context, userId *user_id.UserId) error
}
