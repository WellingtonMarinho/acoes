package alerts

import "time"

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
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	Symbol      string       `json:"symbol"`
	TargetPrice float64      `json:"target_price"`
	Direction   Direction    `json:"direction"`
	DeviceToken string       `json:"device_token,omitempty"`
	Status      AlertStatus  `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	TriggeredAt *time.Time   `json:"triggered_at,omitempty"`
}

type PriceSnapshot struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
