package main

import (
	"context"
	"log"
	"notificationservice/internal/database"
	"notificationservice/internal/kafka"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalln(err)
	}

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	topicDLQ := os.Getenv("KAFKA_TOPIC_DLQ")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	clientID := os.Getenv("KAFKA_CLIENT_ID")
	kafkaRetryes := 3

	dlq := kafka.NewWriter(brokers, topicDLQ, clientID)

	consumer := kafka.NewConsumer(db, dlq, kafkaRetryes, brokers, topic, groupID)
	defer consumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("consuming topic=%s group=%s brokers=%v", topic, groupID, brokers)

	for {
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
