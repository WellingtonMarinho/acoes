package store

import (
	"context"
	"path/filepath"
	"testing"

	"ideacoes/backend/internal/alerts"
)

func TestFileAlertRepositoryPersistsAlerts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alerts.json")

	repo, err := NewFileAlertRepository(path)
	if err != nil {
		t.Fatalf("NewFileAlertRepository() error = %v", err)
	}

	created, err := repo.Create(context.Background(), alerts.Alert{
		ID:          "alert-1",
		UserID:      "user-1",
		Symbol:      "PETR4",
		TargetPrice: 40.5,
		Direction:   alerts.DirectionAbove,
		Status:      alerts.AlertStatusOpen,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != "alert-1" {
		t.Fatalf("expected id alert-1, got %q", created.ID)
	}

	loaded, err := NewFileAlertRepository(path)
	if err != nil {
		t.Fatalf("reload error = %v", err)
	}

	items, err := loaded.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Symbol != "PETR4" {
		t.Fatalf("expected symbol PETR4, got %q", items[0].Symbol)
	}
}
