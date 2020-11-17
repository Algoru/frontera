package rest

import "github.com/Algoru/frontera/domain/service"

// Port
type Port interface {
	Start() error
	SetUserService(service.UserService)
	InitRoutes()
}
