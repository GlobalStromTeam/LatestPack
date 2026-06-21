package main

import (
	"log"
	"os"
	"path/filepath"

	"latestpack/config"
	"latestpack/database"
	"latestpack/repository"
	"latestpack/routes"
	"latestpack/seed"
)

func main() {
	cfg := config.LoadConfig()

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	dbPath := filepath.Join(cfg.DataDir, "latestpack.db")
	db, err := database.Connect(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repos := repository.NewRepositories(db)

	if err := seed.SeedAdminUser(repos.User); err != nil {
		log.Printf("Warning: seed admin user: %v", err)
	}
	if err := seed.SeedLocalChannel(repos.Channel); err != nil {
		log.Printf("Warning: seed local channel: %v", err)
	}

	r := routes.SetupRouter(cfg, repos)

	log.Printf("Server starting on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
