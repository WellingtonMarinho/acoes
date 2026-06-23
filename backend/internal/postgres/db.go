package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	var pingErr error
	for attempts := 0; attempts < 20; attempts++ {
		if pingErr = db.PingContext(ctx); pingErr == nil {
			return db, nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	_ = db.Close()
	return nil, fmt.Errorf("ping postgres: %w", pingErr)
}
