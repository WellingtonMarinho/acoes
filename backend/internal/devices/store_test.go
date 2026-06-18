package devices

import (
	"context"
	"path/filepath"
	"testing"
)

func TestFileRepositoryPersistsRegistrations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")

	repo, err := NewFileRepository(path)
	if err != nil {
		t.Fatalf("NewFileRepository() error = %v", err)
	}

	created, err := repo.Upsert(context.Background(), Registration{
		UserID:      "user-1",
		DeviceToken: "token-1",
		Platform:    "android",
	})
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if created.UserID != "user-1" {
		t.Fatalf("expected user-1, got %q", created.UserID)
	}

	loaded, err := NewFileRepository(path)
	if err != nil {
		t.Fatalf("reload error = %v", err)
	}

	registration, ok, err := loaded.Resolve(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !ok {
		t.Fatalf("expected registration to exist")
	}
	if registration.DeviceToken != "token-1" {
		t.Fatalf("expected token-1, got %q", registration.DeviceToken)
	}
}
