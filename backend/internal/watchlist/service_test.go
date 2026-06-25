package watchlist

import (
	"context"
	"testing"
	"time"

	"ideacoes/backend/internal/actions"
	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/pricefeed"
)

type testActionRepo struct {
	action actions.Action
}

func (r testActionRepo) GetAction(ctx context.Context, id string) (actions.Action, error) {
	_ = ctx
	if r.action.ID == "" || r.action.ID != id {
		return actions.Action{}, actions.ErrActionNotFound
	}
	return r.action, nil
}

type testWatchlistRepo struct {
	items map[string]Item
}

func newTestWatchlistRepo() *testWatchlistRepo {
	return &testWatchlistRepo{items: make(map[string]Item)}
}

func (r *testWatchlistRepo) Upsert(ctx context.Context, item Item) (Item, error) {
	_ = ctx
	r.items[item.UserID+":"+item.ActionID] = item
	return item, nil
}

func (r *testWatchlistRepo) ListByUser(ctx context.Context, userID string) ([]Item, error) {
	_ = ctx
	var out []Item
	for _, item := range r.items {
		if item.UserID == userID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *testWatchlistRepo) Delete(ctx context.Context, userID, actionID string) error {
	_ = ctx
	delete(r.items, userID+":"+actionID)
	return nil
}

type testAlertRepo struct {
	alerts []alerts.Alert
}

func (r *testAlertRepo) ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error) {
	_ = ctx
	var out []alerts.Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID {
			out = append(out, alert)
		}
	}
	return out, nil
}

func (r *testAlertRepo) DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error) {
	_ = ctx
	var deleted int64
	var remaining []alerts.Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID && alert.ActionID == actionID {
			deleted++
			continue
		}
		remaining = append(remaining, alert)
	}
	r.alerts = remaining
	return deleted, nil
}

func TestWatchlistAddListAndRemove(t *testing.T) {
	repo := newTestWatchlistRepo()
	actionsRepo := testActionRepo{action: actions.Action{
		ID:       "action-petr4",
		Symbol:   "PETR4",
		Name:     "Petrobras PN",
		Exchange: "B3",
		Active:   true,
	}}
	alertRepo := &testAlertRepo{
		alerts: []alerts.Alert{
			{UserID: "user-1", ActionID: "action-petr4", Status: alerts.AlertStatusOpen},
			{UserID: "user-1", ActionID: "action-petr4", Status: alerts.AlertStatusTriggered},
		},
	}
	feed := pricefeed.NewMemoryFeed()
	if err := feed.Upsert(context.Background(), alerts.PriceSnapshot{Symbol: "PETR4", Price: 41.2, ObservedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	service := NewService(repo, actionsRepo, alertRepo, feed)

	if _, err := service.Add(context.Background(), "user-1", "action-petr4"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	items, err := service.List(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ActionID != "action-petr4" || items[0].CurrentPrice == nil || *items[0].CurrentPrice != 41.2 {
		t.Fatalf("unexpected item %#v", items[0])
	}
	if items[0].OpenAlertsCount != 1 {
		t.Fatalf("expected 1 open alert, got %d", items[0].OpenAlertsCount)
	}

	if err := service.Remove(context.Background(), "user-1", "action-petr4"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if len(repo.items) != 0 {
		t.Fatalf("expected watchlist to be empty, got %#v", repo.items)
	}
	if len(alertRepo.alerts) != 0 {
		t.Fatalf("expected alerts to be deleted, got %#v", alertRepo.alerts)
	}
}
