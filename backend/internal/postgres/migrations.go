package postgres

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db Execer) error {
	stmts := []string{
		`create table if not exists alerts (
			id text primary key,
			user_id text not null,
			symbol text not null,
			target_price double precision not null,
			direction text not null,
			device_token text not null default '',
			status text not null,
			created_at timestamptz not null,
			triggered_at timestamptz null
		)`,
		`create index if not exists alerts_user_id_idx on alerts (user_id)`,
		`create index if not exists alerts_symbol_status_idx on alerts (symbol, status)`,
		`create table if not exists device_registrations (
			user_id text primary key,
			device_token text not null,
			platform text not null default '',
			created_at timestamptz not null
		)`,
	}

	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
