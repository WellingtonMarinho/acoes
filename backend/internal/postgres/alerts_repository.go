package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	query := `insert into alerts (id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.UserID, alert.ActionID, alert.Symbol, alert.ActionName, alert.TargetPrice, string(alert.Direction),
		alert.DeviceToken, string(alert.Status), alert.CreatedAt, alert.UpdatedAt, alert.TriggeredAt,
	)
	if err != nil {
		return alerts.Alert{}, err
	}
	return alert, nil
}

func (r *AlertRepository) List(ctx context.Context) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts order by created_at asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts where user_id=$1 order by created_at asc`, strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) ListOpenBySymbol(ctx context.Context, symbol string) ([]alerts.Alert, error) {
	rows, err := r.db.QueryContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts where symbol=$1 and status=$2 order by created_at asc`, strings.ToUpper(strings.TrimSpace(symbol)), string(alerts.AlertStatusOpen))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAlerts(rows)
}

func (r *AlertRepository) Get(ctx context.Context, id string) (alerts.Alert, error) {
	row := r.db.QueryRowContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts where id=$1`, strings.TrimSpace(id))
	alert, err := scanAlert(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return alerts.Alert{}, alerts.ErrAlertNotFound
		}
		return alerts.Alert{}, err
	}
	return alert, nil
}

func (r *AlertRepository) Update(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	_, err := r.db.ExecContext(ctx, `update alerts set target_price=$2, direction=$3, updated_at=$4 where id=$1`, alert.ID, alert.TargetPrice, string(alert.Direction), alert.UpdatedAt)
	if err != nil {
		return alerts.Alert{}, err
	}
	return r.Get(ctx, alert.ID)
}

func (r *AlertRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `delete from alerts where id=$1`, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return alerts.ErrAlertNotFound
	}
	return nil
}

func (r *AlertRepository) DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `delete from alerts where user_id=$1 and action_id=$2`, strings.TrimSpace(userID), strings.TrimSpace(actionID))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *AlertRepository) MarkTriggered(ctx context.Context, id string, triggeredAt time.Time) (alerts.Alert, error) {
	result, err := r.db.ExecContext(ctx, `update alerts set status=$2, triggered_at=$3, updated_at=$4 where id=$1 and status=$5`, id, string(alerts.AlertStatusTriggered), triggeredAt, triggeredAt, string(alerts.AlertStatusOpen))
	if err != nil {
		return alerts.Alert{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return alerts.Alert{}, err
	}
	if rows == 0 {
		row := r.db.QueryRowContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts where id=$1`, id)
		alert, err := scanAlert(row)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return alerts.Alert{}, alerts.ErrAlertNotFound
			}
			return alerts.Alert{}, err
		}
		if alert.Status != alerts.AlertStatusOpen {
			return alerts.Alert{}, alerts.ErrAlertNotEditable
		}
		return alerts.Alert{}, alerts.ErrAlertNotFound
	}
	row := r.db.QueryRowContext(ctx, `select id, user_id, action_id, symbol, action_name, target_price, direction, device_token, status, created_at, updated_at, triggered_at from alerts where id=$1`, id)
	alert, err := scanAlert(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return alerts.Alert{}, alerts.ErrAlertNotFound
		}
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
	var actionID string
	var actionName string
	var direction string
	var status string
	var triggeredAt sql.NullTime
	if err := row.Scan(&alert.ID, &alert.UserID, &actionID, &alert.Symbol, &actionName, &alert.TargetPrice, &direction, &alert.DeviceToken, &status, &alert.CreatedAt, &alert.UpdatedAt, &triggeredAt); err != nil {
		return alerts.Alert{}, err
	}
	alert.ActionID = actionID
	alert.ActionName = actionName
	alert.Direction = alerts.Direction(direction)
	alert.Status = alerts.AlertStatus(status)
	if triggeredAt.Valid {
		t := triggeredAt.Time
		alert.TriggeredAt = &t
	}
	return alert, nil
}

var _ alerts.Repository = (*AlertRepository)(nil)
