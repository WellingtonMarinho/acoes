package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/httpapi"
	"ideacoes/backend/internal/memory"
	"ideacoes/backend/internal/monitor"
	"ideacoes/backend/internal/pricefeed"
	"ideacoes/backend/internal/store"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	repo, err := alertRepositoryFromEnv(logger)
	if err != nil {
		logger.Fatalf("repository init error: %v", err)
	}
	deviceRepo, err := deviceRepositoryFromEnv(logger)
	if err != nil {
		logger.Fatalf("device repository init error: %v", err)
	}
	deviceService := devices.NewService(deviceRepo)
	feed := pricefeed.NewMemoryFeed()
	notifier := alerts.NewLogNotifier(logger)
	service := alerts.NewService(repo, notifier, deviceService)
	server := httpapi.NewServer(service, deviceService, feed, logger)
	worker := monitor.NewWorker(service, feed, logger, monitorInterval())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Run(ctx)

	addr := envOrDefault("HTTP_ADDR", ":8080")
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("api listening on %s", addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server error: %v", err)
	}
}

func monitorInterval() time.Duration {
	raw := os.Getenv("MONITOR_INTERVAL_SECONDS")
	if raw == "" {
		return 10 * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return 10 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func alertRepositoryFromEnv(logger *log.Logger) (alerts.Repository, error) {
	if path := os.Getenv("ALERTS_STORE_PATH"); path != "" {
		logger.Printf("using file-backed alert store at %s", path)
		return store.NewFileAlertRepository(path)
	}

	logger.Printf("using in-memory alert store")
	return memory.NewAlertRepository(), nil
}

func deviceRepositoryFromEnv(logger *log.Logger) (devices.Repository, error) {
	if path := os.Getenv("DEVICES_STORE_PATH"); path != "" {
		logger.Printf("using file-backed device store at %s", path)
		return devices.NewFileRepository(path)
	}

	logger.Printf("using in-memory device store")
	return devices.NewMemoryRepository(), nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
