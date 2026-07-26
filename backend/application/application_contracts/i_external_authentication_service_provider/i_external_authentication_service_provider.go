package i_external_authentication_service_provider

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/auth_service_access_token_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/id_token_claims_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/refresh_access_token_response_dto"
)

type IExternalAuthenticationServiceProvider interface {
	GenerateOidcUrl(state string) (string, error)
	GetAccessToken(ctx context.Context, code string) (*auth_service_access_token_response_dto.AuthServiceAccessTokenResponseDto, error)
	ParseIdToken(idToken string) (*id_token_claims_dto.IdTokenClaimsDto, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*refresh_access_token_response_dto.RefreshAccessTokenResponseDto, error)
}
