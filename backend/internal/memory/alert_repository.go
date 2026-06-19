package memory

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/alerts"
)

var errAlertNotFound = errors.New("alert not found")

type AlertRepository struct {
	mu     sync.RWMutex
	alerts map[string]alerts.Alert
}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{
		alerts: make(map[string]alerts.Alert),
	}
}

func (r *AlertRepository) Create(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.alerts[alert.ID] = alert
	return alert, nil
}

func (r *AlertRepository) List(ctx context.Context) ([]alerts.Alert, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]alerts.Alert, 0, len(r.alerts))
	for _, alert := range r.alerts {
		out = append(out, alert)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *AlertRepository) ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error) {
	_ = ctx
	userID = strings.TrimSpace(userID)
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []alerts.Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID {
			out = append(out, alert)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *AlertRepository) ListOpenBySymbol(ctx context.Context, symbol string) ([]alerts.Alert, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	var out []alerts.Alert
	for _, alert := range r.alerts {
		if alert.Symbol == symbol && alert.Status == alerts.AlertStatusOpen {
			out = append(out, alert)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *AlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	alert, ok := r.alerts[id]
	if !ok {
		return alerts.Alert{}, errAlertNotFound
	}

	alert.Status = alerts.AlertStatusTriggered
	alert.TriggeredAt = &triggeredAt
	r.alerts[id] = alert
	return alert, nil
}
