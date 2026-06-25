package postgres

import (
	"context"
	"database/sql"
	"strings"

	"ideacoes/backend/internal/watchlist"
)

type WatchlistRepository struct {
	db *sql.DB
}

func NewWatchlistRepository(db *sql.DB) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

func (r *WatchlistRepository) Upsert(ctx context.Context, item watchlist.Item) (watchlist.Item, error) {
	item.UserID = strings.TrimSpace(item.UserID)
	item.ActionID = strings.TrimSpace(item.ActionID)
	row := r.db.QueryRowContext(ctx, `
		insert into watchlist_items (user_id, action_id, created_at)
		values ($1, $2, coalesce(nullif($3, '0001-01-01 00:00:00+00'::timestamptz), now()))
		on conflict (user_id, action_id) do update set action_id = excluded.action_id
		returning user_id, action_id, created_at
	`, item.UserID, item.ActionID, item.CreatedAt)

	var stored watchlist.Item
	if err := row.Scan(&stored.UserID, &stored.ActionID, &stored.CreatedAt); err != nil {
		return watchlist.Item{}, err
	}
	return stored, nil
}

func (r *WatchlistRepository) ListByUser(ctx context.Context, userID string) ([]watchlist.Item, error) {
	rows, err := r.db.QueryContext(ctx, `select user_id, action_id, created_at from watchlist_items where user_id=$1 order by created_at asc`, strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []watchlist.Item
	for rows.Next() {
		var item watchlist.Item
		if err := rows.Scan(&item.UserID, &item.ActionID, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *WatchlistRepository) Delete(ctx context.Context, userID, actionID string) error {
	result, err := r.db.ExecContext(ctx, `delete from watchlist_items where user_id=$1 and action_id=$2`, strings.TrimSpace(userID), strings.TrimSpace(actionID))
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return watchlist.ErrWatchlistItemNotFound
	}
	return nil
}

var _ watchlist.Repository = (*WatchlistRepository)(nil)
