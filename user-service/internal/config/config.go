package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 5*time.Second),
	}
}

func getDuration(key string, defaultDuration time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultDuration
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultDuration
	}

	return duration
}
