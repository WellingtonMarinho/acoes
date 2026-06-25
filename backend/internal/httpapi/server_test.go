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

	"ideacoes/backend/internal/actions"
	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/memory"
	"ideacoes/backend/internal/pricefeed"
	"ideacoes/backend/internal/watchlist"
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

func (r *alertRepo) Get(ctx context.Context, id string) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.alerts {
		if item.ID == id {
			return item, nil
		}
	}
	return alerts.Alert{}, alerts.ErrAlertNotFound
}

func (r *alertRepo) Update(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.alerts {
		if item.ID == alert.ID {
			r.alerts[i] = alert
			return alert, nil
		}
	}
	return alerts.Alert{}, alerts.ErrAlertNotFound
}

func (r *alertRepo) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.alerts {
		if item.ID == id {
			r.alerts = append(r.alerts[:i], r.alerts[i+1:]...)
			return nil
		}
	}
	return alerts.ErrAlertNotFound
}

func (r *alertRepo) DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var deleted int64
	var remaining []alerts.Alert
	for _, item := range r.alerts {
		if item.UserID == userID && item.ActionID == actionID {
			deleted++
			continue
		}
		remaining = append(remaining, item)
	}
	r.alerts = remaining
	return deleted, nil
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

type watchlistRepo struct {
	mu    sync.Mutex
	items []watchlist.Item
}

func (r *watchlistRepo) Upsert(ctx context.Context, item watchlist.Item) (watchlist.Item, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, existing := range r.items {
		if existing.UserID == item.UserID && existing.ActionID == item.ActionID {
			item.CreatedAt = existing.CreatedAt
			r.items[i] = item
			return item, nil
		}
	}
	r.items = append(r.items, item)
	return item, nil
}

func (r *watchlistRepo) ListByUser(ctx context.Context, userID string) ([]watchlist.Item, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []watchlist.Item
	for _, item := range r.items {
		if item.UserID == userID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *watchlistRepo) Delete(ctx context.Context, userID, actionID string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.items {
		if item.UserID == userID && item.ActionID == actionID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return watchlist.ErrWatchlistItemNotFound
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

func TestCreateAlertIsPublic(t *testing.T) {
	actionService := actions.NewService(memory.NewActionRepository())
	alertRepo := &alertRepo{}
	watchlistService := watchlist.NewService(&watchlistRepo{}, actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, nil, nil, watchlistService, actionService)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, actionService, watchlistService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "")

	req := httptest.NewRequest(http.MethodPost, "/alerts", strings.NewReader(`{"action_id":"action-petr4","target_price":40.5,"direction":"above"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestCreateAlertUsesDefaultUserlessFlow(t *testing.T) {
	alertRepo := &alertRepo{}
	actionService := actions.NewService(memory.NewActionRepository())
	watchlistRepo := &watchlistRepo{}
	watchlistService := watchlist.NewService(watchlistRepo, actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, nil, nil, watchlistService, actionService)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, actionService, watchlistService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "")

	req := httptest.NewRequest(http.MethodPost, "/alerts", strings.NewReader(`{"action_id":"action-petr4","target_price":40.5,"direction":"above"}`))
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
	if alertRepo.alerts[0].UserID != "user-001" {
		t.Fatalf("expected user-001, got %q", alertRepo.alerts[0].UserID)
	}
	if alertRepo.alerts[0].ActionID != "action-petr4" {
		t.Fatalf("expected action-petr4, got %q", alertRepo.alerts[0].ActionID)
	}
}

func TestListAlertsReturnsItems(t *testing.T) {
	alertRepo := &alertRepo{}
	alertRepo.alerts = []alerts.Alert{
		{ID: "1", UserID: "user-001", Symbol: "PETR4"},
		{ID: "2", UserID: "user-999", Symbol: "VALE3"},
	}
	actionService := actions.NewService(memory.NewActionRepository())
	watchlistService := watchlist.NewService(&watchlistRepo{}, actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, nil, nil, watchlistService, actionService)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, actionService, watchlistService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "")

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rec := httptest.NewRecorder()

	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var items []alerts.Alert
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(items) != 1 || items[0].UserID != "user-001" {
		t.Fatalf("expected only alerts for user-001, got %#v", items)
	}
}

func TestEndToEndMVPFlow(t *testing.T) {
	alertRepo := memory.NewAlertRepository()
	deviceRepo := devices.NewMemoryRepository()
	actionService := actions.NewService(memory.NewActionRepository())
	watchlistService := watchlist.NewService(memory.NewWatchlistRepository(), actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, deviceTokenResolver{repo: deviceRepo}, nil, watchlistService, actionService)
	deviceService := devices.NewService(deviceRepo)
	feed := pricefeed.NewMemoryFeed()
	server := NewServer(alertService, actionService, watchlistService, deviceService, feed, log.New(os.Stdout, "", 0), "")

	registerDeviceDirect(t, server.Routes(), "device-123", "android")
	createAlertDirect(t, server.Routes(), "action-petr4", 40.5, "above")
	upsertPriceDirect(t, server.Routes(), "PETR4", 41)

	req, err := http.NewRequest(http.MethodPost, "/prices/check", strings.NewReader(`{"prices":[{"symbol":"PETR4","price":41}]}`))
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

func TestActionsSearchAndWatchlistFlow(t *testing.T) {
	alertRepo := memory.NewAlertRepository()
	actionService := actions.NewService(memory.NewActionRepository())
	watchlistService := watchlist.NewService(memory.NewWatchlistRepository(), actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, nil, nil, watchlistService, actionService)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, actionService, watchlistService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "")

	req := httptest.NewRequest(http.MethodGet, "/actions?query=Petrobras%20PN", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var actionsPayload map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&actionsPayload); err != nil {
		t.Fatalf("decode actions response: %v", err)
	}
	items, ok := actionsPayload["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected one action, got %#v", actionsPayload["items"])
	}

	addReq := httptest.NewRequest(http.MethodPost, "/watchlist", strings.NewReader(`{"action_id":"action-petr4"}`))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(addRec, addReq)
	if addRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", addRec.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
	listRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listRec.Code)
	}
	var watchlistPayload map[string]any
	if err := json.NewDecoder(listRec.Body).Decode(&watchlistPayload); err != nil {
		t.Fatalf("decode watchlist response: %v", err)
	}
	watchlistItems, ok := watchlistPayload["items"].([]any)
	if !ok || len(watchlistItems) != 1 {
		t.Fatalf("expected one watchlist item, got %#v", watchlistPayload["items"])
	}
}

func TestCreateActionAndAddToWatchlist(t *testing.T) {
	actionService := actions.NewService(memory.NewActionRepository())
	alertRepo := memory.NewAlertRepository()
	watchlistService := watchlist.NewService(memory.NewWatchlistRepository(), actionService, alertRepo, pricefeed.NewMemoryFeed())
	alertService := alerts.NewServiceWithActionResolver(alertRepo, noopNotifier{}, nil, nil, watchlistService, actionService)
	deviceService := devices.NewService(&deviceRepo{})
	server := NewServer(alertService, actionService, watchlistService, deviceService, pricefeed.NewMemoryFeed(), log.New(os.Stdout, "", 0), "")

	req := httptest.NewRequest(http.MethodPost, "/actions", strings.NewReader(`{"symbol":"ABCD3","name":"Acao Teste","exchange":"B3"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var created actions.Action
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode created action: %v", err)
	}
	if created.ID == "" || created.Symbol != "ABCD3" {
		t.Fatalf("unexpected action payload: %#v", created)
	}

	addReq := httptest.NewRequest(http.MethodPost, "/watchlist", strings.NewReader(`{"action_id":"`+created.ID+`"}`))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(addRec, addReq)
	if addRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 when adding new action, got %d", addRec.Code)
	}
}

func registerDeviceDirect(t *testing.T, handler http.Handler, deviceToken, platform string) {
	t.Helper()
	body := strings.NewReader(`{"device_token":"` + deviceToken + `","platform":"` + platform + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/devices/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func createAlertDirect(t *testing.T, handler http.Handler, actionID string, targetPrice float64, direction string) {
	t.Helper()
	body := strings.NewReader(`{"action_id":"` + actionID + `","target_price":` + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", targetPrice), "0"), ".") + `,"direction":"` + direction + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/alerts", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func upsertPriceDirect(t *testing.T, handler http.Handler, symbol string, price float64) {
	t.Helper()
	body := strings.NewReader(`{"symbol":"` + symbol + `","price":` + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", price), "0"), ".") + `}`)
	req := httptest.NewRequest(http.MethodPut, "/prices", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func mustRegisterDevice(t *testing.T, client *http.Client, baseURL, deviceToken, platform string) {
	t.Helper()

	body := strings.NewReader(`{"device_token":"` + deviceToken + `","platform":"` + platform + `"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/devices/register", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("register device: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
}

func mustCreateAlert(t *testing.T, client *http.Client, baseURL, actionID string, targetPrice float64, direction string) {
	t.Helper()

	body := strings.NewReader(`{"action_id":"` + actionID + `","target_price":` + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", targetPrice), "0"), ".") + `,"direction":"` + direction + `"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/alerts", body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
