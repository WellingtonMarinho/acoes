package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/alerts"
)

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

func (r *AlertRepository) Get(ctx context.Context, id string) (alerts.Alert, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	alert, ok := r.alerts[id]
	if !ok {
		return alerts.Alert{}, alerts.ErrAlertNotFound
	}
	return alert, nil
}

func (r *AlertRepository) Update(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.alerts[alert.ID]; !ok {
		return alerts.Alert{}, alerts.ErrAlertNotFound
	}
	r.alerts[alert.ID] = alert
	return alert, nil
}

func (r *AlertRepository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.alerts[id]; !ok {
		return alerts.ErrAlertNotFound
	}
	delete(r.alerts, id)
	return nil
}

func (r *AlertRepository) DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	userID = strings.TrimSpace(userID)
	actionID = strings.TrimSpace(actionID)
	var deleted int64
	for id, alert := range r.alerts {
		if alert.UserID == userID && alert.ActionID == actionID {
			delete(r.alerts, id)
			deleted++
		}
	}
	return deleted, nil
}

func (r *AlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	alert, ok := r.alerts[id]
	if !ok {
		return alerts.Alert{}, alerts.ErrAlertNotFound
	}
	if alert.Status != alerts.AlertStatusOpen {
		return alerts.Alert{}, alerts.ErrAlertNotEditable
	}

	alert.Status = alerts.AlertStatusTriggered
	alert.TriggeredAt = &triggeredAt
	alert.UpdatedAt = triggeredAt
	r.alerts[id] = alert
	return alert, nil
}
