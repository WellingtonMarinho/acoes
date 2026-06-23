package postgres

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	"ideacoes/backend/internal/alerts"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	query := `insert into alerts (id, user_id, symbol, target_price, direction, device_token, status, created_at, triggered_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.UserID, alert.Symbol, alert.TargetPrice, string(alert.Direction),
		alert.DeviceToken, string(alert.Status), alert.CreatedAt, alert.TriggeredAt,
	)
	if err != nil {
		return alerts.Alert{}, err
	}
	return alert, nil
}

func (r *AlertRepository) List(ctx context.Context) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, symbol, target_price, direction, device_token, status, created_at, triggered_at from alerts order by created_at asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, symbol, target_price, direction, device_token, status, created_at, triggered_at from alerts where user_id=$1 order by created_at asc`, strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) ListOpenBySymbol(ctx context.Context, symbol string) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, symbol, target_price, direction, device_token, status, created_at, triggered_at from alerts where symbol=$1 and status=$2 order by created_at asc`, strings.ToUpper(strings.TrimSpace(symbol)), string(alerts.AlertStatusOpen))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
	_, err := r.db.ExecContext(ctx, `update alerts set status=$2, triggered_at=$3 where id=$1`, id, string(alerts.AlertStatusTriggered), triggeredAt)
	if err != nil {
		return alerts.Alert{}, err
	}
	row := r.db.QueryRowContext(ctx, `select id, user_id, symbol, target_price, direction, device_token, status, created_at, triggered_at from alerts where id=$1`, id)
	alert, err := scanAlert(row)
	if err != nil {
		return alerts.Alert{}, err
	}
	return alert, nil
}

func scanAlerts(rows *sql.Rows) ([]alerts.Alert, error) {
	out := []alerts.Alert{}
	for rows.Next() {
		alert, err := scanAlert(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, alert)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

type alertRowScanner interface {
	Scan(dest ...any) error
}

func scanAlert(row alertRowScanner) (alerts.Alert, error) {
	var alert alerts.Alert
	var direction string
	var status string
	var triggeredAt sql.NullTime
	if err := row.Scan(&alert.ID, &alert.UserID, &alert.Symbol, &alert.TargetPrice, &direction, &alert.DeviceToken, &status, &alert.CreatedAt, &triggeredAt); err != nil {
		return alerts.Alert{}, err
	}
	alert.Direction = alerts.Direction(direction)
	alert.Status = alerts.AlertStatus(status)
	if triggeredAt.Valid {
		t := triggeredAt.Time
		alert.TriggeredAt = &t
	}
	return alert, nil
}

var _ alerts.Repository = (*AlertRepository)(nil)
