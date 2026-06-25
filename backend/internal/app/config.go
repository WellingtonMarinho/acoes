package app

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr         string
	MonitorInterval  time.Duration
	AlertsStorePath  string
	DevicesStorePath string
	DatabaseURL      string
	PriceFeedProvider string
	TwelveDataAPIKey  string
	TwelveDataBaseURL string
}

func LoadConfig() Config {
	return Config{
		HTTPAddr:         envOrDefault("HTTP_ADDR", ":8080"),
		MonitorInterval:  monitorInterval(),
		AlertsStorePath:  os.Getenv("ALERTS_STORE_PATH"),
		DevicesStorePath: os.Getenv("DEVICES_STORE_PATH"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		PriceFeedProvider: envOrDefault("PRICEFEED_PROVIDER", "memory"),
		TwelveDataAPIKey:  os.Getenv("TWELVEDATA_API_KEY"),
		TwelveDataBaseURL: envOrDefault("TWELVEDATA_BASE_URL", "https://api.twelvedata.com"),
	}
}

func ValidateConfig(cfg Config) error {
	return nil
}

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
