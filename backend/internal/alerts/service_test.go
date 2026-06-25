package alerts

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"ideacoes/backend/internal/actions"
)

type testNotifier struct {
	mu        sync.Mutex
	triggered []Alert
}

type testDeviceResolver struct {
	token string
	ok    bool
}

type testWatchlistRegistrar struct {
	calls []string
}

func (r *testDeviceResolver) Resolve(ctx context.Context, userID string) (string, bool, error) {
	_ = ctx
	_ = userID
	return r.token, r.ok, nil
}

func (r *testWatchlistRegistrar) Upsert(ctx context.Context, userID, actionID string) error {
	_ = ctx
	r.calls = append(r.calls, userID+":"+actionID)
	return nil
}

func (n *testNotifier) Notify(ctx context.Context, alert Alert, marketPrice float64) error {
	_ = ctx
	_ = marketPrice
	n.mu.Lock()
	defer n.mu.Unlock()
	n.triggered = append(n.triggered, alert)
	return nil
}

func TestServiceTriggersAboveAlert(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, nil, &testActionResolver{})

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-petr4",
		TargetPrice: 40,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	triggered, err := service.CheckPrices(context.Background(), []PriceSnapshot{{Symbol: "petr4", Price: 41}})
	if err != nil {
		t.Fatalf("CheckPrices() error = %v", err)
	}
	if len(triggered) != 1 {
		t.Fatalf("expected 1 triggered alert, got %d", len(triggered))
	}
	if triggered[0].ID != created.ID {
		t.Fatalf("expected triggered alert %q, got %q", created.ID, triggered[0].ID)
	}
}

func TestServiceCanTriggerWithoutNotifier(t *testing.T) {
	repo := newTestRepo()
	service := NewServiceWithActionResolver(repo, nil, nil, nil, nil, &testActionResolver{})

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-petr4",
		TargetPrice: 40,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	triggered, err := service.CheckPrices(context.Background(), []PriceSnapshot{{Symbol: "PETR4", Price: 41}})
	if err != nil {
		t.Fatalf("CheckPrices() error = %v", err)
	}
	if len(triggered) != 1 || triggered[0].ID != created.ID {
		t.Fatalf("expected created alert to trigger, got %#v", triggered)
	}
}

func TestServiceDoesNotTriggerBelowAlertEarly(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, nil, &testActionResolver{})

	_, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-vale3",
		TargetPrice: 60,
		Direction:   DirectionBelow,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	triggered, err := service.CheckPrices(context.Background(), []PriceSnapshot{{Symbol: "VALE3", Price: 61}})
	if err != nil {
		t.Fatalf("CheckPrices() error = %v", err)
	}
	if len(triggered) != 0 {
		t.Fatalf("expected no triggered alerts, got %d", len(triggered))
	}
}

func TestServiceIgnoresInvalidPriceSnapshots(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, nil, &testActionResolver{})

	_, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-vale3",
		TargetPrice: 60,
		Direction:   DirectionBelow,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	triggered, err := service.CheckPrices(context.Background(), []PriceSnapshot{{Symbol: "VALE3", Price: 0}})
	if err != nil {
		t.Fatalf("CheckPrices() error = %v", err)
	}
	if len(triggered) != 0 {
		t.Fatalf("expected invalid snapshot to be ignored, got %d triggered alerts", len(triggered))
	}
	if len(notifier.triggered) != 0 {
		t.Fatalf("expected notifier not to be called, got %d calls", len(notifier.triggered))
	}
}

func TestServiceIgnoresAlreadyTriggeredAlertDuringCheck(t *testing.T) {
	repo := newTestRepo()
	repo.markTriggeredErr = ErrAlertNotEditable
	notifier := &testNotifier{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, nil, &testActionResolver{})

	_, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-petr4",
		TargetPrice: 40,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	triggered, err := service.CheckPrices(context.Background(), []PriceSnapshot{{Symbol: "PETR4", Price: 41}})
	if err != nil {
		t.Fatalf("CheckPrices() error = %v", err)
	}
	if len(triggered) != 0 {
		t.Fatalf("expected duplicate trigger to be ignored, got %d", len(triggered))
	}
}

func TestServiceUsesRegisteredDeviceToken(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	resolver := &testDeviceResolver{token: "device-token-123", ok: true}
	service := NewServiceWithActionResolver(repo, notifier, resolver, nil, nil, &testActionResolver{})

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-bbsa3",
		TargetPrice: 12.5,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}
	if created.DeviceToken != "device-token-123" {
		t.Fatalf("expected device token to be filled from resolver, got %q", created.DeviceToken)
	}
}

func TestServiceAddsWatchlistWhenCreatingAlert(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	watchlistRegistrar := &testWatchlistRegistrar{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, watchlistRegistrar, &testActionResolver{})

	_, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-petr4",
		TargetPrice: 40,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}
	if len(watchlistRegistrar.calls) != 1 || watchlistRegistrar.calls[0] != "user-1:action-petr4" {
		t.Fatalf("expected watchlist upsert to be called once, got %#v", watchlistRegistrar.calls)
	}
}

func TestServiceUpdatesAndDeletesAlert(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	service := NewServiceWithActionResolver(repo, notifier, nil, nil, nil, &testActionResolver{})

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		ActionID:    "action-petr4",
		TargetPrice: 40,
		Direction:   DirectionAbove,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	updated, err := service.UpdateAlert(context.Background(), "user-1", created.ID, AlertUpdate{
		TargetPrice: 42,
		Direction:   DirectionBelow,
	})
	if err != nil {
		t.Fatalf("UpdateAlert() error = %v", err)
	}
	if updated.TargetPrice != 42 || updated.Direction != DirectionBelow {
		t.Fatalf("unexpected update result: %#v", updated)
	}

	if err := service.DeleteAlert(context.Background(), "user-1", created.ID); err != nil {
		t.Fatalf("DeleteAlert() error = %v", err)
	}
	if _, err := repo.Get(context.Background(), created.ID); err != ErrAlertNotFound {
		t.Fatalf("expected alert to be deleted, got %v", err)
	}
}

type testRepo struct {
	mu               sync.Mutex
	alerts           map[string]Alert
	markTriggeredErr error
}

func newTestRepo() *testRepo {
	return &testRepo{alerts: make(map[string]Alert)}
}

func (r *testRepo) Create(ctx context.Context, alert Alert) (Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts[alert.ID] = alert
	return alert, nil
}

func (r *testRepo) List(ctx context.Context) ([]Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Alert, 0, len(r.alerts))
	for _, alert := range r.alerts {
		out = append(out, alert)
	}
	return out, nil
}

func (r *testRepo) ListByUser(ctx context.Context, userID string) ([]Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID {
			out = append(out, alert)
		}
	}
	return out, nil
}

func (r *testRepo) ListOpenBySymbol(ctx context.Context, symbol string) ([]Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []Alert
	for _, alert := range r.alerts {
		if alert.Symbol == symbol && alert.Status == AlertStatusOpen {
			out = append(out, alert)
		}
	}
	return out, nil
}

func (r *testRepo) Get(ctx context.Context, id string) (Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	alert, ok := r.alerts[id]
	if !ok {
		return Alert{}, ErrAlertNotFound
	}
	return alert, nil
}

func (r *testRepo) Update(ctx context.Context, alert Alert) (Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.alerts[alert.ID]; !ok {
		return Alert{}, ErrAlertNotFound
	}
	r.alerts[alert.ID] = alert
	return alert, nil
}

func (r *testRepo) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.alerts[id]; !ok {
		return ErrAlertNotFound
	}
	delete(r.alerts, id)
	return nil
}

func (r *testRepo) DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var deleted int64
	for id, alert := range r.alerts {
		if alert.UserID == userID && alert.ActionID == actionID {
			delete(r.alerts, id)
			deleted++
		}
	}
	return deleted, nil
}

func (r *testRepo) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.markTriggeredErr != nil {
		return Alert{}, r.markTriggeredErr
	}
	alert := r.alerts[id]
	if alert.Status != AlertStatusOpen {
		return Alert{}, ErrAlertNotEditable
	}
	alert.Status = AlertStatusTriggered
	alert.TriggeredAt = &triggeredAt
	alert.UpdatedAt = triggeredAt
	r.alerts[id] = alert
	return alert, nil
}

type testActionResolver struct{}

func (r *testActionResolver) GetAction(ctx context.Context, id string) (actions.Action, error) {
	_ = ctx
	return actions.Action{ID: id, Symbol: strings.ToUpper(strings.TrimPrefix(id, "action-")), Name: id, Active: true}, nil
}
