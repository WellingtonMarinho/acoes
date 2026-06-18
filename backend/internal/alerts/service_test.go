package alerts

import (
	"context"
	"sync"
	"testing"
	"time"
)

type testNotifier struct {
	mu        sync.Mutex
	triggered []Alert
}

type testDeviceResolver struct {
	token string
	ok    bool
}

func (r *testDeviceResolver) Resolve(ctx context.Context, userID string) (string, bool, error) {
	_ = ctx
	_ = userID
	return r.token, r.ok, nil
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
	service := NewService(repo, notifier, nil)

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		Symbol:      "PETR4",
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

func TestServiceDoesNotTriggerBelowAlertEarly(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	service := NewService(repo, notifier, nil)

	_, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		Symbol:      "VALE3",
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

func TestServiceUsesRegisteredDeviceToken(t *testing.T) {
	repo := newTestRepo()
	notifier := &testNotifier{}
	resolver := &testDeviceResolver{token: "device-token-123", ok: true}
	service := NewService(repo, notifier, resolver)

	created, err := service.CreateAlert(context.Background(), Alert{
		UserID:      "user-1",
		Symbol:      "B3SA3",
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

type testRepo struct {
	mu     sync.Mutex
	alerts map[string]Alert
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

func (r *testRepo) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	alert := r.alerts[id]
	alert.Status = AlertStatusTriggered
	alert.TriggeredAt = &triggeredAt
	r.alerts[id] = alert
	return alert, nil
}
