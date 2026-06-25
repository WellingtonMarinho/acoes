package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/watchlist"
)

type WatchlistRepository struct {
	mu    sync.RWMutex
	items map[string]map[string]watchlist.Item
}

func NewWatchlistRepository() *WatchlistRepository {
	return &WatchlistRepository{items: make(map[string]map[string]watchlist.Item)}
}

func (r *WatchlistRepository) Upsert(ctx context.Context, item watchlist.Item) (watchlist.Item, error) {
	_ = ctx
	item.UserID = strings.TrimSpace(item.UserID)
	item.ActionID = strings.TrimSpace(item.ActionID)
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	userItems, ok := r.items[item.UserID]
	if !ok {
		userItems = make(map[string]watchlist.Item)
		r.items[item.UserID] = userItems
	}
	if existing, ok := userItems[item.ActionID]; ok {
		item.CreatedAt = existing.CreatedAt
	}
	userItems[item.ActionID] = item
	return item, nil
}

func (r *WatchlistRepository) ListByUser(ctx context.Context, userID string) ([]watchlist.Item, error) {
	_ = ctx
	userID = strings.TrimSpace(userID)

	r.mu.RLock()
	defer r.mu.RUnlock()

	userItems := r.items[userID]
	out := make([]watchlist.Item, 0, len(userItems))
	for _, item := range userItems {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *WatchlistRepository) Delete(ctx context.Context, userID, actionID string) error {
	_ = ctx
	userID = strings.TrimSpace(userID)
	actionID = strings.TrimSpace(actionID)

	r.mu.Lock()
	defer r.mu.Unlock()

	userItems, ok := r.items[userID]
	if !ok {
		return watchlist.ErrWatchlistItemNotFound
	}
	if _, ok := userItems[actionID]; !ok {
		return watchlist.ErrWatchlistItemNotFound
	}
	delete(userItems, actionID)
	if len(userItems) == 0 {
		delete(r.items, userID)
	}
	return nil
}

var _ watchlist.Repository = (*WatchlistRepository)(nil)
