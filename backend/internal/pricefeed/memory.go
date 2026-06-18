package pricefeed

import (
	"context"
	"sort"
	"strings"
	"sync"

	"ideacoes/backend/internal/alerts"
)

type Feed interface {
	List(ctx context.Context) ([]alerts.PriceSnapshot, error)
	Upsert(ctx context.Context, snapshot alerts.PriceSnapshot) error
}

type MemoryFeed struct {
	mu        sync.RWMutex
	snapshots map[string]alerts.PriceSnapshot
}

func NewMemoryFeed() *MemoryFeed {
	return &MemoryFeed{
		snapshots: make(map[string]alerts.PriceSnapshot),
	}
}

func (f *MemoryFeed) List(ctx context.Context) ([]alerts.PriceSnapshot, error) {
	_ = ctx
	f.mu.RLock()
	defer f.mu.RUnlock()

	out := make([]alerts.PriceSnapshot, 0, len(f.snapshots))
	for _, snapshot := range f.snapshots {
		out = append(out, snapshot)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToUpper(out[i].Symbol) < strings.ToUpper(out[j].Symbol)
	})
	return out, nil
}

func (f *MemoryFeed) Upsert(ctx context.Context, snapshot alerts.PriceSnapshot) error {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()

	snapshot.Symbol = strings.ToUpper(strings.TrimSpace(snapshot.Symbol))
	f.snapshots[snapshot.Symbol] = snapshot
	return nil
}
