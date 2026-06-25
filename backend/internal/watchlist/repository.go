package watchlist

import "context"

type Repository interface {
	Upsert(ctx context.Context, item Item) (Item, error)
	ListByUser(ctx context.Context, userID string) ([]Item, error)
	Delete(ctx context.Context, userID, actionID string) error
}
