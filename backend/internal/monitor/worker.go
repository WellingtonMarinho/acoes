package monitor

import (
	"context"
	"log"
	"time"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/pricefeed"
)

type Worker struct {
	service  *alerts.Service
	feed     pricefeed.Feed
	logger   *log.Logger
	interval time.Duration
}

func NewWorker(service *alerts.Service, feed pricefeed.Feed, logger *log.Logger, interval time.Duration) *Worker {
	return &Worker{
		service:  service,
		feed:     feed,
		logger:   logger,
		interval: interval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Printf("price monitor started interval=%s", w.interval)
	w.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			w.logger.Printf("price monitor stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	snapshots, err := w.feed.List(ctx)
	if err != nil {
		w.logger.Printf("price monitor list error: %v", err)
		return
	}
	if len(snapshots) == 0 {
		return
	}

	triggered, err := w.service.CheckPrices(ctx, snapshots)
	if err != nil {
		w.logger.Printf("price monitor check error: %v", err)
		return
	}
	if len(triggered) > 0 {
		w.logger.Printf("price monitor triggered=%d", len(triggered))
	}
}
