package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/alerts"
)

var errAlertNotFound = errors.New("alert not found")

type FileAlertRepository struct {
	mu     sync.RWMutex
	path   string
	alerts map[string]alerts.Alert
}

func NewFileAlertRepository(path string) (*FileAlertRepository, error) {
	repo := &FileAlertRepository{
		path:   path,
		alerts: make(map[string]alerts.Alert),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *FileAlertRepository) Create(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.alerts[alert.ID] = alert
	if err := r.persistLocked(); err != nil {
		return alerts.Alert{}, err
	}
	return alert, nil
}

func (r *FileAlertRepository) List(ctx context.Context) ([]alerts.Alert, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	return sortedAlerts(r.alerts), nil
}

func (r *FileAlertRepository) ListOpenBySymbol(ctx context.Context, symbol string) ([]alerts.Alert, error) {
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

func (r *FileAlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
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
	if err := r.persistLocked(); err != nil {
		return alerts.Alert{}, err
	}
	return alert, nil
}

func (r *FileAlertRepository) load() error {
	if r.path == "" {
		return nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read alerts store: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	var stored []alerts.Alert
	if err := json.Unmarshal(data, &stored); err != nil {
		return fmt.Errorf("decode alerts store: %w", err)
	}

	for _, alert := range stored {
		r.alerts[alert.ID] = alert
	}
	return nil
}

func (r *FileAlertRepository) persistLocked() error {
	if r.path == "" {
		return nil
	}

	stored := sortedAlerts(r.alerts)
	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return fmt.Errorf("encode alerts store: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("create alerts store dir: %w", err)
	}

	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write alerts store tmp: %w", err)
	}
	if err := os.Rename(tmp, r.path); err != nil {
		return fmt.Errorf("replace alerts store: %w", err)
	}
	return nil
}

func sortedAlerts(src map[string]alerts.Alert) []alerts.Alert {
	out := make([]alerts.Alert, 0, len(src))
	for _, alert := range src {
		out = append(out, alert)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out
}
