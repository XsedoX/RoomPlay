package dependecies

import (
	"xsedox.com/main/application"
	"xsedox.com/main/application/user"
	"xsedox.com/main/presentation/handlers"
)

type OidcDependencies struct {
	oidcHandler *handlers.OidcHandler
}

func NewOidcDependencies(loginCommandHandler application.ICommandHandler[*user.LoginCommand]) *OidcDependencies {
	return &OidcDependencies{
		oidcHandler: handlers.NewOidcHandler(loginCommandHandler),
	}
}

func (deps *OidcDependencies) GetOidcHandler() *handlers.OidcHandler {
	return deps.oidcHandler
}
