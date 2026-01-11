package contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos"
)

type IGoogleOidcService interface {
	GenerateOidcUrl(state string) (string, error)
	GetAccessToken(ctx context.Context, code string) (*dtos.GoogleTokenResponseDto, error)
	ParseIdToken(idToken string) (*dtos.GoogleIdTokenClaimsDto, error)
}
