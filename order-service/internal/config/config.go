package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPShutdownTimeout time.Duration
	GRPCRequestTimeout  time.Duration
	KafkaBrokers        []string
	KafkaTopic          string
	KafkaClientID       string
}

func Load() Config {
	return Config{
		HTTPShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 5*time.Second),
		GRPCRequestTimeout:  getDuration("GRPC_REQUEST_TIMEOUT", 2*time.Second),
		KafkaBrokers:        strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:          os.Getenv("KAFKA_ORDERS_CREATED_TOPIC"),
		KafkaClientID:       os.Getenv("KAFKA_CLIENT_ID"),
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
