package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/actions"
)

type ActionRepository struct {
	mu      sync.RWMutex
	actions map[string]actions.Action
}

func NewActionRepository() *ActionRepository {
	now := time.Now().UTC()
	repo := &ActionRepository{
		actions: map[string]actions.Action{
			"action-petr4": {ID: "action-petr4", Symbol: "PETR4", Name: "Petrobras PN", Exchange: "B3", Active: true, CreatedAt: now, UpdatedAt: now},
			"action-vale3": {ID: "action-vale3", Symbol: "VALE3", Name: "Vale ON", Exchange: "B3", Active: true, CreatedAt: now, UpdatedAt: now},
			"action-bbsa3": {ID: "action-bbsa3", Symbol: "BBAS3", Name: "Banco do Brasil ON", Exchange: "B3", Active: true, CreatedAt: now, UpdatedAt: now},
		},
	}
	return repo
}

func (r *ActionRepository) List(ctx context.Context, query string) ([]actions.Action, error) {
	_ = ctx
	query = strings.TrimSpace(query)
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]actions.Action, 0, len(r.actions))
	for _, item := range r.actions {
		if !item.Active {
			continue
		}
		if query != "" && !strings.EqualFold(item.Name, query) {
			continue
		}
		if item.Active {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Symbol < out[j].Symbol
	})
	return out, nil
}

func (r *ActionRepository) Get(ctx context.Context, id string) (actions.Action, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.actions[strings.TrimSpace(id)]
	if !ok {
		return actions.Action{}, actions.ErrActionNotFound
	}
	return item, nil
}

func (r *ActionRepository) Upsert(ctx context.Context, action actions.Action) (actions.Action, error) {
	_ = ctx
	action.Symbol = strings.ToUpper(strings.TrimSpace(action.Symbol))
	action.Name = strings.TrimSpace(action.Name)
	action.Exchange = strings.TrimSpace(action.Exchange)
	if action.Symbol == "" || action.Name == "" {
		return actions.Action{}, actions.ErrInvalidAction
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	if existing, ok := r.actions[strings.TrimSpace(action.ID)]; ok {
		action.CreatedAt = existing.CreatedAt
	} else {
		for _, existing := range r.actions {
			if strings.EqualFold(existing.Symbol, action.Symbol) {
				action.ID = existing.ID
				action.CreatedAt = existing.CreatedAt
				break
			}
		}
	}
	if action.ID == "" {
		action.ID = "action-" + strings.ToLower(action.Symbol)
	}
	action.Active = true
	action.UpdatedAt = now
	if action.CreatedAt.IsZero() {
		action.CreatedAt = now
	}
	r.actions[action.ID] = action
	return action, nil
}

var _ actions.Repository = (*ActionRepository)(nil)
