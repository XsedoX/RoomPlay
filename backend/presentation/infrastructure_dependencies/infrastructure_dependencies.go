package infrastructure_dependencies

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_google_oidc_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/google_oidc_service"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/jwt_provider"
)

type InfrastructureDependencies struct {
	Encrypter         i_encrypter.IEncrypter
	JwtProvider       i_jwt_provider.IJwtProvider
	GoogleOidcService i_google_oidc_service.IGoogleOidcService
}

func ConstructInfrastructureDependencies(config config.IConfiguration) *InfrastructureDependencies {
	encrypter := encryper.NewEncrypter(config)
	googleOidcService := google_oidc_service.NewGoogleOidcService(config)
	jwtProvider := jwt_provider.NewJwtProvider(config)

	return &InfrastructureDependencies{
		Encrypter:         encrypter,
		GoogleOidcService: googleOidcService,
		JwtProvider:       jwtProvider,
	}
}
