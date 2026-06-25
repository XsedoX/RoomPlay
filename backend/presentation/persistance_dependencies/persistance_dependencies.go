package persistance_dependencies

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/application/user/user_contracts/i_user_repository"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/external_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/room/room_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_repository"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/jmoiron/sqlx"
)

type PersistanceDependencies struct {
	ExternalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
	InternalCredentialsRepository i_internal_credentials_repository.IInternalCredentialsRepository
	UserRepository                i_user_repository.IUserRepository
	RoomRepository                i_room_repository.IRoomRepository
	UnitOfWork                    i_unit_of_work.IUnitOfWork
}

func ConstructPersistanceDependencies(
	ctx context.Context,
	db *sqlx.DB,
	configuration config.IConfiguration,
	infraDependencies *infrastructure_dependencies.InfrastructureDependencies,
) *PersistanceDependencies {
	unitOfWork := unit_of_work.NewUnitOfWork(db)

	userRepository := user_repository.NewUserRepository()
	externalCredentialsRepository := external_credentials_repository.NewExternalCredentialsRepository(infraDependencies.Encrypter)
	internalCredentialsRepository := internal_credentials_repository.NewInternalCredentialsRepository(infraDependencies.Encrypter)
	roomRepository := room_repository.NewRoomRepository(infraDependencies.Encrypter)

	return &PersistanceDependencies{
		ExternalCredentialsRepository: externalCredentialsRepository,
		InternalCredentialsRepository: internalCredentialsRepository,
		UserRepository:                userRepository,
		RoomRepository:                roomRepository,
		UnitOfWork:                    unitOfWork,
	}
}
