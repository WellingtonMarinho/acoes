package postgres

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	"ideacoes/backend/internal/devices"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Upsert(ctx context.Context, registration devices.Registration) (devices.Registration, error) {
	registration.UserID = strings.TrimSpace(registration.UserID)
	registration.DeviceToken = strings.TrimSpace(registration.DeviceToken)
	registration.Platform = strings.TrimSpace(registration.Platform)
	if registration.UserID == "" || registration.DeviceToken == "" {
		return devices.Registration{}, devices.ErrInvalidDeviceRegistration
	}
	if registration.CreatedAt.IsZero() {
		registration.CreatedAt = sqlNow()
	}

	query := `insert into device_registrations (user_id, device_token, platform, created_at)
		values ($1,$2,$3,$4)
		on conflict (user_id) do update set device_token = excluded.device_token, platform = excluded.platform, created_at = excluded.created_at`
	_, err := r.db.ExecContext(ctx, query, registration.UserID, registration.DeviceToken, registration.Platform, registration.CreatedAt)
	if err != nil {
		return devices.Registration{}, err
	}
	return registration, nil
}

func (r *DeviceRepository) Resolve(ctx context.Context, userID string) (devices.Registration, bool, error) {
	row := r.db.QueryRowContext(ctx, `select user_id, device_token, platform, created_at from device_registrations where user_id=$1`, strings.TrimSpace(userID))
	var registration devices.Registration
	if err := row.Scan(&registration.UserID, &registration.DeviceToken, &registration.Platform, &registration.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return devices.Registration{}, false, nil
		}
		return devices.Registration{}, false, err
	}
	return registration, true, nil
}

func (r *DeviceRepository) List(ctx context.Context) ([]devices.Registration, error) {
	rows, err := r.db.QueryContext(ctx, `select user_id, device_token, platform, created_at from device_registrations order by created_at asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRegistrations(rows)
}

func (r *DeviceRepository) ListByUser(ctx context.Context, userID string) ([]devices.Registration, error) {
	rows, err := r.db.QueryContext(ctx, `select user_id, device_token, platform, created_at from device_registrations where user_id=$1 order by created_at asc`, strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRegistrations(rows)
}

func scanRegistrations(rows *sql.Rows) ([]devices.Registration, error) {
	out := []devices.Registration{}
	for rows.Next() {
		var registration devices.Registration
		if err := rows.Scan(&registration.UserID, &registration.DeviceToken, &registration.Platform, &registration.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, registration)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func sqlNow() time.Time {
	return time.Now().UTC()
}

var _ devices.Repository = (*DeviceRepository)(nil)
