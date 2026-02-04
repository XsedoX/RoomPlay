package i_google_oidc_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/google_id_token_claims_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/google_token_response_dto"
)

type IGoogleOidcService interface {
	GenerateOidcUrl(state string) (string, error)
	GetAccessToken(ctx context.Context, code string) (*google_token_response_dto.GoogleTokenResponseDto, error)
	ParseIdToken(idToken string) (*google_id_token_claims_dto.GoogleIdTokenClaimsDto, error)
}
