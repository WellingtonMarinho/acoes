package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ideacoes/backend/internal/actions"
)

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) List(ctx context.Context, query string) ([]actions.Action, error) {
	query = strings.TrimSpace(query)
	sqlQuery := `select id, symbol, name, exchange, active, created_at, updated_at from actions where active = true`
	args := []any{}
	if query != "" {
		sqlQuery += ` and lower(name) = lower($1)`
		args = append(args, query)
	}
	sqlQuery += ` order by symbol asc`

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []actions.Action
	for rows.Next() {
		var item actions.Action
		if err := rows.Scan(&item.ID, &item.Symbol, &item.Name, &item.Exchange, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ActionRepository) Get(ctx context.Context, id string) (actions.Action, error) {
	var item actions.Action
	err := r.db.QueryRowContext(ctx, `select id, symbol, name, exchange, active, created_at, updated_at from actions where id = $1`, strings.TrimSpace(id)).
		Scan(&item.ID, &item.Symbol, &item.Name, &item.Exchange, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return actions.Action{}, actions.ErrActionNotFound
		}
		return actions.Action{}, err
	}
	return item, nil
}

func (r *ActionRepository) Upsert(ctx context.Context, action actions.Action) (actions.Action, error) {
	action.Symbol = strings.ToUpper(strings.TrimSpace(action.Symbol))
	action.Name = strings.TrimSpace(action.Name)
	action.Exchange = strings.TrimSpace(action.Exchange)
	if action.Symbol == "" || action.Name == "" {
		return actions.Action{}, actions.ErrInvalidAction
	}

	if strings.TrimSpace(action.ID) == "" {
		action.ID = "action-" + strings.ToLower(action.Symbol)
	}

	row := r.db.QueryRowContext(ctx, `
		insert into actions (id, symbol, name, exchange, active, created_at, updated_at)
		values ($1, $2, $3, $4, true, now(), now())
		on conflict (symbol) do update set
			name = excluded.name,
			exchange = excluded.exchange,
			active = true,
			updated_at = now()
		returning id, symbol, name, exchange, active, created_at, updated_at`,
		action.ID, action.Symbol, action.Name, action.Exchange)

	var stored actions.Action
	if err := row.Scan(&stored.ID, &stored.Symbol, &stored.Name, &stored.Exchange, &stored.Active, &stored.CreatedAt, &stored.UpdatedAt); err != nil {
		return actions.Action{}, err
	}
	return stored, nil
}

var _ actions.Repository = (*ActionRepository)(nil)
