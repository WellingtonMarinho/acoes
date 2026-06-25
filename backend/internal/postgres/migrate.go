package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	return MigrateFromDir(ctx, db, migrationDir())
}

func MigrateFromDir(ctx context.Context, db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.UpContext(ctx, db, dir)
}

func migrationDir() string {
	if dir := os.Getenv("GOOSE_DIR"); dir != "" {
		return dir
	}

	if wd, err := os.Getwd(); err == nil {
		candidates := []string{
			filepath.Join(wd, "migrations"),
			filepath.Join(wd, "internal", "postgres", "migrations"),
		}
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	return "internal/postgres/migrations"
}
