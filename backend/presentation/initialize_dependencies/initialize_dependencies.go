package initialize_dependencies

import (
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/presentation_dependencies"
)

type ServerDependencies struct {
	PresentationDependencies   *presentation_dependencies.PresentationDependencies
	InfrastructureDependencies *infrastructure_dependencies.InfrastructureDependencies
}

func NewServerDependencies(
	infraDependencies *infrastructure_dependencies.InfrastructureDependencies,
	presentationDependencies *presentation_dependencies.PresentationDependencies,
) *ServerDependencies {
	return &ServerDependencies{
		PresentationDependencies:   presentationDependencies,
		InfrastructureDependencies: infraDependencies,
	}
}
