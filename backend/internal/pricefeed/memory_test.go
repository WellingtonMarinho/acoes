package pricefeed

import (
	"context"
	"testing"

	"ideacoes/backend/internal/alerts"
)

func TestMemoryFeedRegisterSymbolDoesNotCreateSnapshot(t *testing.T) {
	feed := NewMemoryFeed()

	if err := feed.RegisterSymbol(context.Background(), "PETR4"); err != nil {
		t.Fatalf("RegisterSymbol() error = %v", err)
	}

	snapshots, err := feed.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(snapshots) != 0 {
		t.Fatalf("expected no snapshots before an observed price, got %#v", snapshots)
	}
}

func TestMemoryFeedListsObservedSnapshots(t *testing.T) {
	feed := NewMemoryFeed()

	if err := feed.Upsert(context.Background(), alerts.PriceSnapshot{Symbol: "petr4", Price: 41.2}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	snapshots, err := feed.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("expected one snapshot, got %d", len(snapshots))
	}
	if snapshots[0].Symbol != "PETR4" || snapshots[0].Price != 41.2 {
		t.Fatalf("unexpected snapshot %#v", snapshots[0])
	}
}
