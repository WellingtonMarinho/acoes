package main

import (
	"context"
	"log"
	"os"

	"ideacoes/backend/internal/app"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	cfg := app.LoadConfig()
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatalf("app init error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := application.Run(ctx); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
