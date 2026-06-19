package app

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr         string
	JWTSecret        string
	MonitorInterval  time.Duration
	AlertsStorePath  string
	DevicesStorePath string
}

func LoadConfig() Config {
	return Config{
		HTTPAddr:         envOrDefault("HTTP_ADDR", ":8080"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		MonitorInterval:  monitorInterval(),
		AlertsStorePath:  os.Getenv("ALERTS_STORE_PATH"),
		DevicesStorePath: os.Getenv("DEVICES_STORE_PATH"),
	}
}

func ValidateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return ErrMissingJWTSecret
	}
	return nil
}

var ErrMissingJWTSecret = errors.New("JWT_SECRET is required")

func monitorInterval() time.Duration {
	raw := os.Getenv("MONITOR_INTERVAL_SECONDS")
	if raw == "" {
		return 10 * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return 10 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
