package main

import (
	"log"
	"os"
	"strconv"

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

	strHTTPRestPort := os.Getenv("REST_HTTP_PORT")
	httpRestPort, err := strconv.ParseUint(strHTTPRestPort, 10, 16)
	if err != nil {
		log.Fatalf("unable to use %s as HTTP REST port: %s\n", strHTTPRestPort, err)
	}

	ginAdapter := ginadapter.NewGinAdapter(uint16(httpRestPort))

	appConfig := application.Configuration{
		DatabasePort: &mongoAdapter,
		RestPort:     &ginAdapter,
	}

	if err := appConfig.Start(); err != nil {
		log.Fatal(err)
	}
}
