package actions

import "context"

type Repository interface {
	List(ctx context.Context, query string) ([]Action, error)
	Get(ctx context.Context, id string) (Action, error)
	Upsert(ctx context.Context, action Action) (Action, error)
}
