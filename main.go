package main

import (
	"log"
	"os"

	"github.com/Algoru/frontera/application"
	mongoadapter "github.com/Algoru/frontera/infrastructure/database/mongo_adapter"
	ginadapter "github.com/Algoru/frontera/infrastructure/rest/gin_adapter"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoAdapter := mongoadapter.MongoAdapter{
		User:     os.Getenv("MONGO_DATABASE_USER"),
		Password: os.Getenv("MONGO_DATABASE_PASSWORD"),
		Host:     os.Getenv("MONGO_DATABASE_HOST"),
		Database: os.Getenv("MONGO_DATABASE_NAME"),
	}

	ginAdapter := ginadapter.GinAdapter{}

	appConfig := application.ApplicationConfiguration{
		DatabaseAdapter: mongoAdapter,
		RestAdatapter:   ginAdapter,
	}

	if err := appConfig.Start(); err != nil {
		log.Fatal(err)
	}
}
