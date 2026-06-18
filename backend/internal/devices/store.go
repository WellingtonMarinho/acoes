package devices

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
)

type FileRepository struct {
	mu    sync.RWMutex
	path  string
	items map[string]Registration
}

func NewFileRepository(path string) (*FileRepository, error) {
	repo := &FileRepository{
		path:  path,
		items: make(map[string]Registration),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *FileRepository) Upsert(ctx context.Context, registration Registration) (Registration, error) {
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
	if err := r.persistLocked(); err != nil {
		return Registration{}, err
	}
	return registration, nil
}

func (r *FileRepository) Resolve(ctx context.Context, userID string) (Registration, bool, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	registration, ok := r.items[strings.TrimSpace(userID)]
	return registration, ok, nil
}

func (r *FileRepository) List(ctx context.Context) ([]Registration, error) {
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

func (r *FileRepository) load() error {
	if r.path == "" {
		return nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read devices store: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	var stored []Registration
	if err := json.Unmarshal(data, &stored); err != nil {
		return fmt.Errorf("decode devices store: %w", err)
	}

	for _, item := range stored {
		r.items[item.UserID] = item
	}
	return nil
}

func (r *FileRepository) persistLocked() error {
	if r.path == "" {
		return nil
	}

	stored := make([]Registration, 0, len(r.items))
	for _, item := range r.items {
		stored = append(stored, item)
	}
	sort.Slice(stored, func(i, j int) bool {
		return stored[i].CreatedAt.Before(stored[j].CreatedAt)
	})

	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return fmt.Errorf("encode devices store: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("create devices store dir: %w", err)
	}

	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write devices store tmp: %w", err)
	}
	if err := os.Rename(tmp, r.path); err != nil {
		return fmt.Errorf("replace devices store: %w", err)
	}
	return nil
}
