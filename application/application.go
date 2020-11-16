package application

import (
	"github.com/Algoru/frontera/infrastructure/database"
	"github.com/Algoru/frontera/infrastructure/rest"
)

// ApplicationConfiguration
type ApplicationConfiguration struct {
	DatabaseAdapter database.DatabasePort
	RestAdatapter   rest.RestPort
}

// Start
func (ac *ApplicationConfiguration) Start() error {
	if err := ac.DatabaseAdapter.StartDatabase(); err != nil {
		return err
	}

	return nil
}
