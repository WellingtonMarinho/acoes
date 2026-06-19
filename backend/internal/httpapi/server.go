package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/auth"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/pricefeed"
)

type Server struct {
	service   *alerts.Service
	devices   *devices.Service
	feed      pricefeed.Feed
	logger    *log.Logger
	jwtSecret string
}

func NewServer(service *alerts.Service, devices *devices.Service, feed pricefeed.Feed, logger *log.Logger, jwtSecret string) *Server {
	return &Server{service: service, devices: devices, feed: feed, logger: logger, jwtSecret: jwtSecret}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /auth/token", s.handleIssueToken)
	mux.Handle("GET /alerts", s.requireAuth(http.HandlerFunc(s.handleListAlerts)))
	mux.Handle("POST /alerts", s.requireAuth(http.HandlerFunc(s.handleCreateAlert)))
	mux.Handle("GET /devices", s.requireAuth(http.HandlerFunc(s.handleListDevices)))
	mux.Handle("POST /devices/register", s.requireAuth(http.HandlerFunc(s.handleRegisterDevice)))
	mux.HandleFunc("GET /prices", s.handleListPrices)
	mux.HandleFunc("PUT /prices", s.handleUpsertPrice)
	mux.HandleFunc("POST /prices/check", s.handleCheckPrices)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	alertsList, err := s.service.ListAlertsByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, alertsList)
}

func (s *Server) handleCreateAlert(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symbol      string           `json:"symbol"`
		TargetPrice float64          `json:"target_price"`
		Direction   alerts.Direction `json:"direction"`
		DeviceToken string           `json:"device_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	userID, _ := auth.UserIDFromContext(r.Context())
	created, err := s.service.CreateAlert(r.Context(), alerts.Alert{
		UserID:      userID,
		Symbol:      req.Symbol,
		TargetPrice: req.TargetPrice,
		Direction:   req.Direction,
		DeviceToken: req.DeviceToken,
	})
	if err != nil {
		status := http.StatusBadRequest
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleListDevices(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
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
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	userID, _ := auth.UserIDFromContext(r.Context())
	created, err := s.devices.Register(r.Context(), devices.Registration{
		UserID:      userID,
		DeviceToken: req.DeviceToken,
		Platform:    req.Platform,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleCheckPrices(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prices []alerts.PriceSnapshot `json:"prices"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
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
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Symbol == "" || req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "invalid price snapshot")
		return
	}
	if err := s.feed.Upsert(r.Context(), req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (s *Server) handleIssueToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, err := auth.Sign(req.UserID, s.jwtSecret, 24*time.Hour, time.Now())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.BearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization")
			return
		}

		claims, err := auth.Parse(token, s.jwtSecret, time.Now())
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := auth.WithUserID(r.Context(), claims.UserID)
		ctx = auth.WithClaims(ctx, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
