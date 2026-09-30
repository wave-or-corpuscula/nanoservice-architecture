package main

import (
	"context"
	"log"
	"notificationservice/internal/kafka"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	consumer := kafka.NewConsumer(brokers, topic, groupID)
	defer consumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		consumer := kafka.NewConsumer(brokers, topic, groupID)
		log.Printf("consuming topic=%s group=%s brokers=%v", topic, groupID, brokers)
		err := consumer.Run(ctx)
		consumer.Close()

		if ctx.Err() != nil {
			log.Println("consumer stopped gracefully")
			return
		}

		log.Printf("consumer error: %v, retrying in 5s...", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}
