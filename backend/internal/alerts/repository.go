package alerts

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, alert Alert) (Alert, error)
	List(ctx context.Context) ([]Alert, error)
	ListByUser(ctx context.Context, userID string) ([]Alert, error)
	ListOpenBySymbol(ctx context.Context, symbol string) ([]Alert, error)
	Get(ctx context.Context, id string) (Alert, error)
	Update(ctx context.Context, alert Alert) (Alert, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error)
	MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (Alert, error)
}
