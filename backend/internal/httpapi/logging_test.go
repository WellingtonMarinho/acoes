package httpapi

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddlewareLogsRequestSummary(t *testing.T) {
	var out bytes.Buffer
	logger := log.New(&out, "", 0)
	handler := loggingMiddleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/alerts?status=open", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	line := out.String()
	for _, want := range []string{
		"event=http_request",
		"method=POST",
		"path=/alerts?status=open",
		"status=201",
		"bytes=2",
		"duration=",
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("expected log to contain %q, got %q", want, line)
		}
	}
}

func TestLoggingMiddlewareDefaultsStatusToOK(t *testing.T) {
	var out bytes.Buffer
	logger := log.New(&out, "", 0)
	handler := loggingMiddleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !strings.Contains(out.String(), "status=200") {
		t.Fatalf("expected default status 200 in log, got %q", out.String())
	}
}
