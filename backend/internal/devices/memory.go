package devices

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidDeviceRegistration = errors.New("invalid device registration")

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]Registration
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: make(map[string]Registration)}
}

func (r *MemoryRepository) Upsert(ctx context.Context, registration Registration) (Registration, error) {
	_ = ctx

	registration.UserID = strings.TrimSpace(registration.UserID)
	registration.DeviceToken = strings.TrimSpace(registration.DeviceToken)
	registration.Platform = strings.TrimSpace(registration.Platform)
	if registration.UserID == "" || registration.DeviceToken == "" {
		return Registration{}, ErrInvalidDeviceRegistration
	}
	if registration.CreatedAt.IsZero() {
		registration.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[registration.UserID] = registration
	return registration, nil
}

func (r *MemoryRepository) Resolve(ctx context.Context, userID string) (Registration, bool, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	registration, ok := r.items[strings.TrimSpace(userID)]
	return registration, ok, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Registration, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Registration, 0, len(r.items))
	for _, item := range r.items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}
