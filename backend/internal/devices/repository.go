package devices

import "context"

type Repository interface {
	Upsert(ctx context.Context, registration Registration) (Registration, error)
	Resolve(ctx context.Context, userID string) (Registration, bool, error)
	List(ctx context.Context) ([]Registration, error)
}
