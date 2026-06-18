package devices

import "time"

type Registration struct {
	UserID      string    `json:"user_id"`
	DeviceToken string    `json:"device_token"`
	Platform    string    `json:"platform,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
