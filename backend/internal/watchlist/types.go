package watchlist

import "time"

type Item struct {
	UserID    string    `json:"user_id"`
	ActionID  string    `json:"action_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ActionID        string     `json:"action_id"`
	Symbol          string     `json:"symbol"`
	Name            string     `json:"name"`
	Exchange        string     `json:"exchange"`
	CurrentPrice    *float64   `json:"current_price,omitempty"`
	LastPriceAt     *time.Time `json:"last_price_at,omitempty"`
	OpenAlertsCount int        `json:"open_alerts_count"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
}
