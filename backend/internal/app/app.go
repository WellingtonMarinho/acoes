package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"ideacoes/backend/internal/actions"
	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/httpapi"
	"ideacoes/backend/internal/memory"
	"ideacoes/backend/internal/monitor"
	backendpostgres "ideacoes/backend/internal/postgres"
	"ideacoes/backend/internal/pricefeed"
	"ideacoes/backend/internal/store"
	"ideacoes/backend/internal/watchlist"
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

	db, err := databaseFromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}

	alertRepo, err := alertRepositoryFromConfig(cfg, db, logger)
	if err != nil {
		return nil, err
	}
	actionRepo, err := actionRepositoryFromConfig(cfg, db, logger)
	if err != nil {
		return nil, err
	}
	actionService := actions.NewService(actionRepo)
	deviceRepo, err := deviceRepositoryFromConfig(cfg, db, logger)
	if err != nil {
		return nil, err
	}

	deviceService := devices.NewService(deviceRepo)
	feed, err := priceFeedFromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}
	watchlistRepo, err := watchlistRepositoryFromConfig(cfg, db, logger)
	if err != nil {
		return nil, err
	}
	watchlistService := watchlist.NewService(watchlistRepo, actionService, alertRepo, feed)
	notifier := alerts.NewLogNotifier(logger)
	alertService := alerts.NewServiceWithActionResolver(alertRepo, notifier, deviceService, feed, watchlistService, actionService)

	server := httpapi.NewServer(alertService, actionService, watchlistService, deviceService, feed, logger, "")
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

func databaseFromConfig(cfg Config, logger *log.Logger) (*sql.DB, error) {
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return nil, nil
	}

	logger.Printf("using postgres database at %s", cfg.DatabaseURL)
	db, err := backendpostgres.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func alertRepositoryFromConfig(cfg Config, db *sql.DB, logger *log.Logger) (alerts.Repository, error) {
	if db != nil {
		logger.Printf("using postgres alert store")
		return backendpostgres.NewAlertRepository(db), nil
	}
	if cfg.AlertsStorePath != "" {
		logger.Printf("using file-backed alert store at %s", cfg.AlertsStorePath)
		return store.NewFileAlertRepository(cfg.AlertsStorePath)
	}

	logger.Printf("using in-memory alert store")
	return memory.NewAlertRepository(), nil
}

func deviceRepositoryFromConfig(cfg Config, db *sql.DB, logger *log.Logger) (devices.Repository, error) {
	if db != nil {
		logger.Printf("using postgres device store")
		return backendpostgres.NewDeviceRepository(db), nil
	}
	if cfg.DevicesStorePath != "" {
		logger.Printf("using file-backed device store at %s", cfg.DevicesStorePath)
		return devices.NewFileRepository(cfg.DevicesStorePath)
	}

	logger.Printf("using in-memory device store")
	return devices.NewMemoryRepository(), nil
}

func actionRepositoryFromConfig(cfg Config, db *sql.DB, logger *log.Logger) (actions.Repository, error) {
	if db != nil {
		logger.Printf("using postgres action store")
		return backendpostgres.NewActionRepository(db), nil
	}

	logger.Printf("using in-memory action store")
	return memory.NewActionRepository(), nil
}

func watchlistRepositoryFromConfig(cfg Config, db *sql.DB, logger *log.Logger) (watchlist.Repository, error) {
	if db != nil {
		logger.Printf("using postgres watchlist store")
		return backendpostgres.NewWatchlistRepository(db), nil
	}

	logger.Printf("using in-memory watchlist store")
	return memory.NewWatchlistRepository(), nil
}

func priceFeedFromConfig(cfg Config, logger *log.Logger) (pricefeed.Feed, error) {
	switch cfg.PriceFeedProvider {
	case "", "memory":
		logger.Printf("using in-memory price feed")
		return pricefeed.NewMemoryFeed(), nil
	case "twelvedata":
		if strings.TrimSpace(cfg.TwelveDataAPIKey) == "" {
			logger.Printf("twelvedata api key missing, falling back to in-memory price feed")
			return pricefeed.NewMemoryFeed(), nil
		}
		logger.Printf("using twelvedata price feed")
		return pricefeed.NewTwelveDataFeed(http.DefaultClient, cfg.TwelveDataBaseURL, cfg.TwelveDataAPIKey), nil
	default:
		logger.Printf("unknown price feed provider %q, using in-memory price feed", cfg.PriceFeedProvider)
		return pricefeed.NewMemoryFeed(), nil
	}
}
