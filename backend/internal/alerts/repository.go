package alerts

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, alert Alert) (Alert, error)
	List(ctx context.Context) ([]Alert, error)
	ListOpenBySymbol(ctx context.Context, symbol string) ([]Alert, error)
	MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (Alert, error)
}
