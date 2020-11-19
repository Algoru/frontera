package application

import (
	"log"

	"github.com/Algoru/frontera/domain/service"
	"github.com/Algoru/frontera/infrastructure/database"
	"github.com/Algoru/frontera/infrastructure/rest"
)

// Configuration
type Configuration struct {
	DatabasePort database.Port
	RestPort     rest.Port
}

// Start
func (ac *Configuration) Start() error {
	log.Println("starting database adapter")
	if err := ac.DatabasePort.StartDatabase(); err != nil {
		return err
	}

	userService := service.NewUserService(ac.DatabasePort)
	authService := service.NewAuthService(ac.DatabasePort, ac.DatabasePort)

	ac.RestPort.SetUserService(userService)
	ac.RestPort.SetAuthService(authService)
	ac.RestPort.InitRoutes()

	log.Println("starting rest adapter")
	return ac.RestPort.Start()
}
