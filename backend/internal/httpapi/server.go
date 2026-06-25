package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ideacoes/backend/internal/actions"
	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/pricefeed"
	"ideacoes/backend/internal/watchlist"
)

type Server struct {
	service   *alerts.Service
	actions   *actions.Service
	watchlist *watchlist.Service
	devices   *devices.Service
	feed      pricefeed.Feed
	logger    *log.Logger
}

func NewServer(service *alerts.Service, actions *actions.Service, watchlist *watchlist.Service, devices *devices.Service, feed pricefeed.Feed, logger *log.Logger, _ string) *Server {
	return &Server{service: service, actions: actions, watchlist: watchlist, devices: devices, feed: feed, logger: logger}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /actions", s.handleListActions)
	mux.Handle("POST /actions", http.HandlerFunc(s.handleCreateAction))
	mux.Handle("GET /watchlist", http.HandlerFunc(s.handleListWatchlist))
	mux.Handle("POST /watchlist", http.HandlerFunc(s.handleAddWatchlist))
	mux.Handle("DELETE /watchlist/{action_id}", http.HandlerFunc(s.handleDeleteWatchlist))
	mux.Handle("GET /alerts", http.HandlerFunc(s.handleListAlerts))
	mux.Handle("POST /alerts", http.HandlerFunc(s.handleCreateAlert))
	mux.Handle("PATCH /alerts/{id}", http.HandlerFunc(s.handleUpdateAlert))
	mux.Handle("DELETE /alerts/{id}", http.HandlerFunc(s.handleDeleteAlert))
	mux.Handle("GET /devices", http.HandlerFunc(s.handleListDevices))
	mux.Handle("POST /devices/register", http.HandlerFunc(s.handleRegisterDevice))
	mux.HandleFunc("GET /prices", s.handleListPrices)
	mux.HandleFunc("PUT /prices", s.handleUpsertPrice)
	mux.HandleFunc("POST /prices/check", s.handleCheckPrices)
	return loggingMiddleware(s.logger, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	userID := "user-001"
	alertsList, err := s.service.ListAlertsByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, alertsList)
}

func (s *Server) handleListActions(w http.ResponseWriter, r *http.Request) {
	items, err := s.actions.ListActions(r.Context(), r.URL.Query().Get("query"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) handleCreateAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symbol   string `json:"symbol"`
		Name     string `json:"name"`
		Exchange string `json:"exchange"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	created, err := s.actions.CreateAction(r.Context(), req.Symbol, req.Name, req.Exchange)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	s.logger.Printf("event=action_created action_id=%s symbol=%s name=%s exchange=%s", created.ID, created.Symbol, created.Name, created.Exchange)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleListWatchlist(w http.ResponseWriter, r *http.Request) {
	userID := "user-001"
	items, err := s.watchlist.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) handleAddWatchlist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string `json:"action_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	userID := "user-001"
	item, err := s.watchlist.Add(r.Context(), userID, req.ActionID)
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) handleDeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	userID := "user-001"
	actionID := r.PathValue("action_id")
	if err := s.watchlist.Remove(r.Context(), userID, actionID); err != nil {
		s.writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleCreateAlert(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID    string           `json:"action_id"`
		TargetPrice float64          `json:"target_price"`
		Direction   alerts.Direction `json:"direction"`
		DeviceToken string           `json:"device_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	userID := "user-001"
	created, err := s.service.CreateAlert(r.Context(), alerts.Alert{
		ActionID:    req.ActionID,
		UserID:      userID,
		TargetPrice: req.TargetPrice,
		Direction:   req.Direction,
		DeviceToken: req.DeviceToken,
	})
	if err != nil {
		s.logger.Printf("event=alert_create_failed user_id=%s action_id=%s error=%v", userID, req.ActionID, err)
		s.writeDomainError(w, err)
		return
	}

	s.logger.Printf("event=alert_created user_id=%s alert_id=%s action_id=%s symbol=%s target=%.2f direction=%s", userID, created.ID, created.ActionID, created.Symbol, created.TargetPrice, created.Direction)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleUpdateAlert(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TargetPrice float64          `json:"target_price"`
		Direction   alerts.Direction `json:"direction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	userID := "user-001"
	updated, err := s.service.UpdateAlert(r.Context(), userID, r.PathValue("id"), alerts.AlertUpdate{
		TargetPrice: req.TargetPrice,
		Direction:   req.Direction,
	})
	if err != nil {
		s.writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteAlert(w http.ResponseWriter, r *http.Request) {
	userID := "user-001"
	if err := s.service.DeleteAlert(r.Context(), userID, r.PathValue("id")); err != nil {
		s.writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListDevices(w http.ResponseWriter, r *http.Request) {
	userID := "user-001"
	items, err := s.devices.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceToken string `json:"device_token"`
		Platform    string `json:"platform"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	userID := "user-001"
	created, err := s.devices.Register(r.Context(), devices.Registration{
		UserID:      userID,
		DeviceToken: req.DeviceToken,
		Platform:    req.Platform,
	})
	if err != nil {
		s.logger.Printf("event=device_register_failed user_id=%s error=%v", userID, err)
		writeAPIError(w, http.StatusBadRequest, "invalid_device_registration", "Registro de device invalido.")
		return
	}

	s.logger.Printf("event=device_registered user_id=%s platform=%s", userID, created.Platform)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleCheckPrices(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prices []alerts.PriceSnapshot `json:"prices"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	triggered, err := s.service.CheckPrices(r.Context(), req.Prices)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"triggered": triggered,
		"count":     len(triggered),
	})
}

func (s *Server) handleListPrices(w http.ResponseWriter, r *http.Request) {
	prices, err := s.feed.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, prices)
}

func (s *Server) handleUpsertPrice(w http.ResponseWriter, r *http.Request) {
	var req alerts.PriceSnapshot
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Printf("event=price_upsert_invalid_json")
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}
	if req.Symbol == "" || req.Price <= 0 {
		s.logger.Printf("event=price_upsert_rejected symbol=%s price=%.2f", req.Symbol, req.Price)
		writeAPIError(w, http.StatusBadRequest, "invalid_price_snapshot", "invalid price snapshot")
		return
	}
	if err := s.feed.Upsert(r.Context(), req); err != nil {
		s.logger.Printf("event=price_upsert_failed symbol=%s error=%v", req.Symbol, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.logger.Printf("event=price_upserted symbol=%s price=%.2f", req.Symbol, req.Price)
	writeJSON(w, http.StatusOK, req)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeAPIError(w, status, statusCodeFor(status), message)
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func statusCodeFor(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	default:
		return "internal_error"
	}
}

func (s *Server) writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, alerts.ErrInvalidAlert):
		writeAPIError(w, http.StatusBadRequest, "invalid_alert", "alerta invalido.")
	case errors.Is(err, alerts.ErrAlertNotFound):
		writeAPIError(w, http.StatusNotFound, "alert_not_found", "Alerta nao encontrado.")
	case errors.Is(err, alerts.ErrAlertNotEditable):
		writeAPIError(w, http.StatusConflict, "alert_not_editable", "Alerta nao pode ser alterado.")
	case errors.Is(err, watchlist.ErrWatchlistItemNotFound):
		writeAPIError(w, http.StatusNotFound, "watchlist_item_not_found", "Acao nao encontrada na watchlist.")
	case errors.Is(err, watchlist.ErrInvalidWatchlistItem):
		writeAPIError(w, http.StatusBadRequest, "invalid_watchlist_item", "Item da watchlist invalido.")
	case errors.Is(err, actions.ErrInvalidAction):
		writeAPIError(w, http.StatusBadRequest, "invalid_action", "Acao invalida.")
	case errors.Is(err, actions.ErrActionNotFound):
		writeAPIError(w, http.StatusNotFound, "action_not_found", "Acao nao encontrada.")
	default:
		s.logger.Printf("event=domain_error err=%v", err)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "erro interno.")
	}
}
