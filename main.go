package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-authorization-service/config"
	"github.com/tnqbao/gau-authorization-service/controller"
	"github.com/tnqbao/gau-authorization-service/infra"
	"github.com/tnqbao/gau-authorization-service/repository"
	"github.com/tnqbao/gau-authorization-service/routes"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system environment variables")
	}

	// Initialize config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize infrastructure (Redis, Postgres, etc.)
	infra, err := infra.InitInfra(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}

	// Initialize repository
	repo, err := repository.InitRepository(infra)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize controller
	ctrl := controller.NewController(cfg, infra, repo)

	// Setup and run HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s", port)
	if err := routes.SetupRouter(ctrl).Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
