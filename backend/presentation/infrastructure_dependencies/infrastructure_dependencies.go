package infrastructure_dependencies

import (
	"context"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_google_oidc_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/google_oidc_service"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/jwt_provider"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache/cache_cleanup_worker"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache/caching_song_decorator"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/external_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/room/room_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/youtube_music_data_provider"
	"github.com/jmoiron/sqlx"
)

type InfrastructureDependencies struct {
	Encrypter                     i_encrypter.IEncrypter
	JwtProvider                   i_jwt_provider.IJwtProvider
	GoogleOidcService             i_google_oidc_service.IGoogleOidcService
	ExternalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
	InternalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
	UserRepository                i_user_repository.IUserRepository
	RoomRepository                i_room_repository.IRoomRepository
	UnitOfWork                    i_unit_of_work.IUnitOfWork
	CachingSongDecorator          i_music_data_provider_service.IMusicDataProviderService
}

func ConstructInfrastructureDependencies(
	ctx context.Context,
	db *sqlx.DB,
	config config.IConfiguration,
) *InfrastructureDependencies {
	go cache_cleanup_worker.StartCacheCleanupWorker(ctx, db, time.Hour, time.Minute*15)
	encrypter := encryper.NewEncrypter(config)
	googleOidcService := google_oidc_service.NewGoogleOidcService(config)
	jwtProvider := jwt_provider.NewJwtProvider(config)
	unitOfWork := unit_of_work.NewUnitOfWork(db)

	searchSongCache := cache.NewCache[[]music_data_response_dto.MusicDataResponseDto](
		config.CacheSimilarityThreshold(),
	)
	youtubeMusicDataProvider := youtube_music_data_provider.NewYoutubeMusicDataProvider()
	cachingSongDecorator := caching_song_decorator.NewCachingSongDecorator(
		youtubeMusicDataProvider,
		searchSongCache,
		unitOfWork,
	)

	userRepository := user_repository.NewUserRepository()
	externalCredentialsRepository := external_credentials_repository.NewExternalCredentialsRepository(encrypter)
	internalCredentialsRepository := internal_credentials_repository.NewInternalCredentialsRepository(encrypter)
	roomRepository := room_repository.NewRoomRepository(encrypter)

	return &InfrastructureDependencies{
		Encrypter:                     encrypter,
		GoogleOidcService:             googleOidcService,
		ExternalCredentialsRepository: externalCredentialsRepository,
		InternalCredentialsRepository: internalCredentialsRepository,
		UserRepository:                userRepository,
		RoomRepository:                roomRepository,
		UnitOfWork:                    unitOfWork,
		CachingSongDecorator:          cachingSongDecorator,
		JwtProvider:                   jwtProvider,
	}
}
