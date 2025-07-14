package main

import (
	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-authorization-service/config"
	"github.com/tnqbao/gau-authorization-service/controller"
	"github.com/tnqbao/gau-authorization-service/infra"
	"github.com/tnqbao/gau-authorization-service/repository"
	"github.com/tnqbao/gau-authorization-service/routes"
	"log"
)

func main() {
	err := godotenv.Load("/gau_authorization/authorization.env")
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	// Initialize configuration and infrastructure
	newConfig := config.InitConfig()
	if newConfig == nil {
		log.Fatal("Failed to initialize configuration")
	}

	newInfra := infra.InitInfra(newConfig)

	if newInfra == nil {
		log.Fatal("Failed to initialize infrastructure")
	}

	newRepository := repository.InitRepository(newInfra)
	if newRepository == nil {
		log.Fatal("Failed to initialize repository")
	}
	// Initialize controller with the new configuration and infrastructure
	ctrl := controller.NewController(newConfig, newInfra, newRepository)

	router := routes.SetupRouter(ctrl)
	router.Run(":8080")
}
