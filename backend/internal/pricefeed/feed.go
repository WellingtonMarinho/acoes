package pricefeed

import (
	"context"

	"ideacoes/backend/internal/alerts"
)

type Feed interface {
	List(ctx context.Context) ([]alerts.PriceSnapshot, error)
	Upsert(ctx context.Context, snapshot alerts.PriceSnapshot) error
	RegisterSymbol(ctx context.Context, symbol string) error
}
