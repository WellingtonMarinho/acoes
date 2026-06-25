//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	testcontainerspostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/devices"
	"ideacoes/backend/internal/watchlist"
)

func TestAlertRepositoryPersistsAgainstPostgres(t *testing.T) {
	ctx := context.Background()
	db, cleanup := openTestDatabase(t, ctx)
	defer cleanup()

	repo := NewAlertRepository(db)
	watchlistRepo := NewWatchlistRepository(db)
	if _, err := watchlistRepo.Upsert(ctx, watchlist.Item{UserID: "user-1", ActionID: "action-petr4", CreatedAt: time.Now().UTC().Truncate(time.Second)}); err != nil {
		t.Fatalf("watchlist Upsert() error = %v", err)
	}

	alert := alerts.Alert{
		ID:          "alert-1",
		UserID:      "user-1",
		ActionID:    "action-petr4",
		Symbol:      "PETR4",
		TargetPrice: 40.5,
		Direction:   alerts.DirectionAbove,
		Status:      alerts.AlertStatusOpen,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
	}

	created, err := repo.Create(ctx, alert)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != alert.ID {
		t.Fatalf("expected id %q, got %q", alert.ID, created.ID)
	}

	items, err := repo.ListByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("ListByUser() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(items))
	}
	if items[0].Symbol != "PETR4" {
		t.Fatalf("expected PETR4, got %q", items[0].Symbol)
	}

	triggeredAt := time.Now().UTC().Truncate(time.Second)
	updated, err := repo.MarkTriggered(ctx, alert.ID, triggeredAt)
	if err != nil {
		t.Fatalf("MarkTriggered() error = %v", err)
	}
	if updated.Status != alerts.AlertStatusTriggered {
		t.Fatalf("expected triggered status, got %q", updated.Status)
	}
	if updated.TriggeredAt == nil {
		t.Fatal("expected triggered_at to be set")
	}

	updatedAt := updated.UpdatedAt
	if updatedAt.IsZero() {
		t.Fatal("expected updated_at to be set")
	}

	if _, err := repo.MarkTriggered(ctx, alert.ID, time.Now().UTC().Truncate(time.Second)); err != alerts.ErrAlertNotEditable {
		t.Fatalf("expected second trigger attempt to be rejected, got %v", err)
	}

	if err := repo.Delete(ctx, alert.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.Get(ctx, alert.ID); err != alerts.ErrAlertNotFound {
		t.Fatalf("expected alert to be deleted, got %v", err)
	}

	created, err = repo.Create(ctx, alert)
	if err != nil {
		t.Fatalf("Create() after delete error = %v", err)
	}
	if err := watchlistRepo.Delete(ctx, "user-1", "action-petr4"); err != nil {
		t.Fatalf("Delete watchlist() error = %v", err)
	}
	if _, err := repo.Get(ctx, created.ID); err != alerts.ErrAlertNotFound {
		t.Fatalf("expected cascade delete from watchlist, got %v", err)
	}
}

func TestDeviceRepositoryPersistsAgainstPostgres(t *testing.T) {
	ctx := context.Background()
	db, cleanup := openTestDatabase(t, ctx)
	defer cleanup()

	repo := NewDeviceRepository(db)

	created, err := repo.Upsert(ctx, devices.Registration{
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

	loaded, ok, err := repo.Resolve(ctx, "user-1")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !ok {
		t.Fatal("expected registration to exist")
	}
	if loaded.DeviceToken != "token-1" {
		t.Fatalf("expected token-1, got %q", loaded.DeviceToken)
	}
}

func TestWatchlistRepositoryPersistsAgainstPostgres(t *testing.T) {
	ctx := context.Background()
	db, cleanup := openTestDatabase(t, ctx)
	defer cleanup()

	actionRepo := NewActionRepository(db)
	action, err := actionRepo.Get(ctx, "action-petr4")
	if err != nil {
		t.Fatalf("Get action() error = %v", err)
	}

	repo := NewWatchlistRepository(db)
	created, err := repo.Upsert(ctx, watchlist.Item{UserID: "user-1", ActionID: action.ID, CreatedAt: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if created.UserID != "user-1" || created.ActionID != action.ID {
		t.Fatalf("unexpected item %#v", created)
	}

	items, err := repo.ListByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("ListByUser() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 watchlist item, got %d", len(items))
	}

	if err := repo.Delete(ctx, "user-1", action.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	items, err = repo.ListByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("ListByUser() after delete error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected watchlist empty, got %d", len(items))
	}
}

func openTestDatabase(t *testing.T, ctx context.Context) (*sql.DB, func()) {
	t.Helper()

	container, err := testcontainerspostgres.Run(ctx,
		"postgres:16-alpine",
		testcontainerspostgres.WithDatabase("ideacoes"),
		testcontainerspostgres.WithUsername("ideacoes"),
		testcontainerspostgres.WithPassword("ideacoes"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("connection string: %v", err)
	}

	db, err := Open(ctx, dsn)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("open db: %v", err)
	}
	if err := Migrate(ctx, db); err != nil {
		_ = db.Close()
		_ = container.Terminate(ctx)
		t.Fatalf("migrate db: %v", err)
	}

	cleanup := func() {
		_ = db.Close()
		_ = container.Terminate(ctx)
	}
	return db, cleanup
}
