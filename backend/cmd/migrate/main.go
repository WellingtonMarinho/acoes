package main

import (
	"context"
	"log"
	"os"

	"ideacoes/backend/internal/app"
	backendpostgres "ideacoes/backend/internal/postgres"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	cfg := app.LoadConfig()
	if cfg.DatabaseURL == "" {
		logger.Fatal("DATABASE_URL is required for migrations")
	}

	db, err := backendpostgres.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := backendpostgres.Migrate(context.Background(), db); err != nil {
		logger.Fatalf("goose up: %v", err)
	}

	logger.Println("migrations applied successfully")
}
