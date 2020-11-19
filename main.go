package main

import (
	"flag"
	"log"

	"github.com/Algoru/frontera/application"
	"github.com/Algoru/frontera/configuration"
	mongoadapter "github.com/Algoru/frontera/infrastructure/database/mongo_adapter"
	ginadapter "github.com/Algoru/frontera/infrastructure/rest/gin_adapter"
)

func main() {
	configFilePath := flag.String("config", "cfg.toml", "specify custom Frontera configuration file")
	flag.Parse()

	log.Println("parsing configuration file")
	if err := configuration.LoadConfigurationFromFile(*configFilePath); err != nil {
		log.Fatal(err)
	}

	configuration.PrintWarnings()

	dbConfig := configuration.GetConfiguration().Database
	mongoAdapter := mongoadapter.MongoAdapter{
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Host:     dbConfig.Host,
		Database: dbConfig.Database,
	}

	httpPort := configuration.GetConfiguration().HTTP.Port
	ginAdapter := ginadapter.NewGinAdapter(httpPort)

	appConfig := application.Configuration{
		DatabasePort: &mongoAdapter,
		RestPort:     &ginAdapter,
	}

	log.Println("starting application")
	if err := appConfig.Start(); err != nil {
		log.Fatal(err)
	}
}
