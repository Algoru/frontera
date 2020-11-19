package configuration

import (
	"log"
	"os"

	"github.com/naoina/toml"
)

// Configuration
type Configuration struct {
	Debug    bool
	Security SecurityConfiguration
	Database DatabaseConfiguration
	User     UserConfiguration
	HTTP     HTTPConfiguration
}

var global Configuration

// LoadConfigurationFromFile
func LoadConfigurationFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = toml.NewDecoder(file).Decode(&global); err != nil {
		return err
	}

	return nil
}

func PrintWarnings() {
	if global.Security.TokenSigningKey == "" {
		log.Println("WARNING: using empty token signing key")
	}
}

// GetConfiguration returns a copy of the loaded configuration
func GetConfiguration() Configuration {
	return global
}
