package initialize_dependencies

import (
	"context"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/presentation/application_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/persistance_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/presentation_dependencies"
	"github.com/jmoiron/sqlx"
)

type ServerDependencies struct {
	PresentationDependencies   *presentation_dependencies.PresentationDependencies
	PersistanceDependencies    *persistance_dependencies.PersistanceDependencies
	InfrastructureDependencies *infrastructure_dependencies.InfrastructureDependencies
}

func NewServerDependencies(ctx context.Context, db *sqlx.DB, configuration config.IConfiguration) *ServerDependencies {
	infrastructureDependendies := infrastructure_dependencies.ConstructInfrastructureDependencies(configuration)
	persistanceDependencies := persistance_dependencies.ConstructPersistanceDependencies(
		ctx,
		db,
		configuration,
		infrastructureDependendies,
	)
	applicationDependencies := application_dependencies.ConstructApplicationDependencies(
		persistanceDependencies,
		infrastructureDependendies,
		configuration,
	)
	presentationDependencies := presentation_dependencies.ConstructPresentationDependencies(
		configuration,
		applicationDependencies,
		infrastructureDependendies,
	)

	return &ServerDependencies{
		PresentationDependencies:   presentationDependencies,
		PersistanceDependencies:    persistanceDependencies,
		InfrastructureDependencies: infrastructureDependendies,
	}
}
