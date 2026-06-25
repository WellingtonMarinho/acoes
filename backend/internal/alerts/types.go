package alerts

import (
	"errors"
	"time"
)

var (
	ErrInvalidAlert     = errors.New("invalid alert")
	ErrAlertNotFound    = errors.New("alert not found")
	ErrAlertNotEditable = errors.New("alert not editable")
)

type Direction string

const (
	DirectionAbove Direction = "above"
	DirectionBelow Direction = "below"
)

type AlertStatus string

const (
	AlertStatusOpen      AlertStatus = "open"
	AlertStatusTriggered AlertStatus = "triggered"
)

type Alert struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	ActionID    string      `json:"action_id"`
	Symbol      string      `json:"symbol"`
	ActionName  string      `json:"action_name,omitempty"`
	TargetPrice float64     `json:"target_price"`
	Direction   Direction   `json:"direction"`
	DeviceToken string      `json:"device_token,omitempty"`
	Status      AlertStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	TriggeredAt *time.Time  `json:"triggered_at,omitempty"`
}

type AlertUpdate struct {
	TargetPrice float64   `json:"target_price"`
	Direction   Direction `json:"direction"`
}

type PriceSnapshot struct {
	Symbol     string    `json:"symbol"`
	Price      float64   `json:"price"`
	ObservedAt time.Time `json:"observed_at,omitempty"`
}
