package mock_external_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/dtos/refresh_access_token_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockExternalCredentialsRepository struct {
	mock.Mock
}

func (m *MockExternalCredentialsRepository) GetRefreshTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error) {
	args := m.Called(ctx, userId, queryer)
	return args.Get(0).(*token.Token), args.Error(1)
}

func (m *MockExternalCredentialsRepository) Grant(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, credentials, queryer)
	return args.Error(0)
}

func (m *MockExternalCredentialsRepository) GetAccessTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error) {
	args := m.Called(ctx, userId, queryer)
	return args.Get(0).(*token.Token), args.Error(1)
}

func (m *MockExternalCredentialsRepository) RefreshAccessToken(ctx context.Context, dto refresh_access_token_dto.RefreshAccessTokenDto, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, dto, queryer)
	return args.Error(0)
}

func (m *MockExternalCredentialsRepository) GetExternalCredentialsByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*external_credentials.ExternalCredentials, error) {
	args := m.Called(ctx, userId, queryer)
	return args.Get(0).(*external_credentials.ExternalCredentials), args.Error(1)
}

func (m *MockExternalCredentialsRepository) UpdateExternalCredentials(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error {
	args := m.Called(ctx, credentials, queryer)
	return args.Error(0)
}
