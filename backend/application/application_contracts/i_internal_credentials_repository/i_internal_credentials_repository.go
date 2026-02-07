package i_internal_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type IInternalCredentialsRepository interface {
	AssignNewToken(ctx context.Context, refreshToken *internal_credentials.InternalCredentials, queryer i_queryer.IQueryer) error
	GetTokenByValue(ctx context.Context, value string, queryer i_queryer.IQueryer) (*internal_credentials.InternalCredentials, error)
	RetireTokenByUserSession(ctx context.Context, userSession user_session.UserSession, queryer i_queryer.IQueryer) error
	RetireAllTokensByUserId(ctx context.Context, userId *user_id.UserId, queryer i_queryer.IQueryer) error
}
