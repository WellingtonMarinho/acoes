package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/auth"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/pricefeed"
)

type alertRepo struct {
	mu     sync.Mutex
	alerts []alerts.Alert
}

func (r *alertRepo) Create(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts = append(r.alerts, alert)
	return alert, nil
}

func (r *alertRepo) List(ctx context.Context) ([]alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]alerts.Alert(nil), r.alerts...), nil
}

func (r *alertRepo) ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []alerts.Alert
	for _, item := range r.alerts {
		if item.UserID == userID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *alertRepo) ListOpenBySymbol(ctx context.Context, symbol string) ([]alerts.Alert, error) {
	_ = ctx
	_ = symbol
	return nil, nil
}

func (r *alertRepo) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
	_ = ctx
	_ = id
	_ = triggeredAt
	return alerts.Alert{}, nil
}

type deviceRepo struct {
	mu            sync.Mutex
	registrations []devices.Registration
}

func (r *deviceRepo) Upsert(ctx context.Context, registration devices.Registration) (devices.Registration, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registrations = append(r.registrations, registration)
	return registration, nil
}

func (r *deviceRepo) Resolve(ctx context.Context, userID string) (devices.Registration, bool, error) {
	_ = ctx
	_ = userID
	return devices.Registration{}, false, nil
}

func (r *deviceRepo) List(ctx context.Context) ([]devices.Registration, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]devices.Registration(nil), r.registrations...), nil
}

func (r *deviceRepo) ListByUser(ctx context.Context, userID string) ([]devices.Registration, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []devices.Registration
	for _, item := range r.registrations {
		if item.UserID == userID {
			out = append(out, item)
		}
	}
	return out, nil
}

type noopNotifier struct{}

func (noopNotifier) Notify(ctx context.Context, alert alerts.Alert, marketPrice float64) error {
	_ = ctx
	_ = alert
	_ = marketPrice
	return nil
}

func TestCreateAlertRequiresJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	alertService := alerts.NewService(&alertRepo{}, noopNotifier{}, nil)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "secret")

	req := httptest.NewRequest(http.MethodPost, "/alerts", strings.NewReader(`{"symbol":"PETR4","target_price":40.5,"direction":"above"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestCreateAlertUsesAuthenticatedUser(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	alertRepo := &alertRepo{}
	alertService := alerts.NewService(alertRepo, noopNotifier{}, nil)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "secret")

	token, err := auth.Sign("user-123", "secret", time.Hour, time.Now())
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/alerts", strings.NewReader(`{"symbol":"PETR4","target_price":40.5,"direction":"above"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	alertRepo.mu.Lock()
	defer alertRepo.mu.Unlock()
	if len(alertRepo.alerts) != 1 {
		t.Fatalf("expected 1 alert stored, got %d", len(alertRepo.alerts))
	}
	if alertRepo.alerts[0].UserID != "user-123" {
		t.Fatalf("expected user-123, got %q", alertRepo.alerts[0].UserID)
	}
}

func TestIssueToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	alertService := alerts.NewService(&alertRepo{}, noopNotifier{}, nil)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "secret")

	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(`{"user_id":"user-123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["access_token"] == "" {
		t.Fatal("expected access_token in response")
	}
}

func TestListAlertsRequiresJWTAndFiltersByUser(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	alertRepo := &alertRepo{}
	alertRepo.alerts = []alerts.Alert{
		{ID: "1", UserID: "user-123", Symbol: "PETR4"},
		{ID: "2", UserID: "user-999", Symbol: "VALE3"},
	}
	alertService := alerts.NewService(alertRepo, noopNotifier{}, nil)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "secret")

	token, err := auth.Sign("user-123", "secret", time.Hour, time.Now())
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var items []alerts.Alert
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(items) != 1 || items[0].UserID != "user-123" {
		t.Fatalf("expected only alerts for user-123, got %#v", items)
	}
}
