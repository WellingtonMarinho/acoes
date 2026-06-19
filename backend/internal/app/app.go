package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/httpapi"
	"ideacoes/backend/internal/memory"
	"ideacoes/backend/internal/monitor"
	"ideacoes/backend/internal/pricefeed"
	"ideacoes/backend/internal/store"
)

type App struct {
	logger     *log.Logger
	httpServer *http.Server
	worker     *monitor.Worker
}

func New(cfg Config, logger *log.Logger) (*App, error) {
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	alertRepo, err := alertRepositoryFromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}
	deviceRepo, err := deviceRepositoryFromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}

	deviceService := devices.NewService(deviceRepo)
	feed := pricefeed.NewMemoryFeed()
	notifier := alerts.NewLogNotifier(logger)
	alertService := alerts.NewService(alertRepo, notifier, deviceService)

	server := httpapi.NewServer(alertService, deviceService, feed, logger, cfg.JWTSecret)
	worker := monitor.NewWorker(alertService, feed, logger, cfg.MonitorInterval)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		logger:     logger,
		httpServer: httpServer,
		worker:     worker,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go a.worker.Run(ctx)

	a.logger.Printf("api listening on %s", a.httpServer.Addr)
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func alertRepositoryFromConfig(cfg Config, logger *log.Logger) (alerts.Repository, error) {
	if cfg.AlertsStorePath != "" {
		logger.Printf("using file-backed alert store at %s", cfg.AlertsStorePath)
		return store.NewFileAlertRepository(cfg.AlertsStorePath)
	}

	logger.Printf("using in-memory alert store")
	return memory.NewAlertRepository(), nil
}

func deviceRepositoryFromConfig(cfg Config, logger *log.Logger) (devices.Repository, error) {
	if cfg.DevicesStorePath != "" {
		logger.Printf("using file-backed device store at %s", cfg.DevicesStorePath)
		return devices.NewFileRepository(cfg.DevicesStorePath)
	}

	logger.Printf("using in-memory device store")
	return devices.NewMemoryRepository(), nil
}
