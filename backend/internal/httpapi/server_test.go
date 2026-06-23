package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
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
	"ideacoes/backend/internal/memory"
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

type deviceTokenResolver struct {
	repo devices.Repository
}

func (r deviceTokenResolver) Resolve(ctx context.Context, userID string) (string, bool, error) {
	registration, ok, err := r.repo.Resolve(ctx, userID)
	if err != nil || !ok {
		return "", ok, err
	}
	return registration.DeviceToken, true, nil
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

func TestEndToEndMVPFlow(t *testing.T) {
	alertRepo := memory.NewAlertRepository()
	deviceRepo := devices.NewMemoryRepository()
	alertService := alerts.NewService(alertRepo, noopNotifier{}, deviceTokenResolver{repo: deviceRepo})
	deviceService := devices.NewService(deviceRepo)
	feed := pricefeed.NewMemoryFeed()
	server := NewServer(alertService, deviceService, feed, log.New(os.Stdout, "", 0), "secret")

	ts := httptest.NewServer(server.Routes())
	defer ts.Close()

	client := ts.Client()

	token := mustIssueToken(t, client, ts.URL, "user-123")
	mustRegisterDevice(t, client, ts.URL, token, "device-123", "android")
	mustCreateAlert(t, client, ts.URL, token, "PETR4", 40.5, "above")
	mustUpsertPrice(t, client, ts.URL, "PETR4", 41)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/prices/check", strings.NewReader(`{"prices":[{"symbol":"PETR4","price":41}]}`))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["count"] != float64(1) {
		t.Fatalf("expected 1 triggered alert, got %#v", payload["count"])
	}
}

func mustIssueToken(t *testing.T, client *http.Client, baseURL, userID string) string {
	t.Helper()

	reqBody := strings.NewReader(`{"user_id":"` + userID + `"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/auth/token", reqBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var payload map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	if payload["access_token"] == "" {
		t.Fatal("expected access token")
	}
	return payload["access_token"]
}

func mustRegisterDevice(t *testing.T, client *http.Client, baseURL, token, deviceToken, platform string) {
	t.Helper()

	body := strings.NewReader(`{"device_token":"` + deviceToken + `","platform":"` + platform + `"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/devices/register", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("register device: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
}

func mustCreateAlert(t *testing.T, client *http.Client, baseURL, token, symbol string, targetPrice float64, direction string) {
	t.Helper()

	body := strings.NewReader(`{"symbol":"` + symbol + `","target_price":` + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", targetPrice), "0"), ".") + `,"direction":"` + direction + `"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/alerts", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("create alert: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
}

func mustUpsertPrice(t *testing.T, client *http.Client, baseURL, symbol string, price float64) {
	t.Helper()

	body := strings.NewReader(`{"symbol":"` + symbol + `","price":` + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", price), "0"), ".") + `}`)
	req, err := http.NewRequest(http.MethodPut, baseURL+"/prices", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("upsert price: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}
