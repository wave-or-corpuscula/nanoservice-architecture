package kafka

import (
	"context"
	"encoding/json"
	"orderservice/internal/database"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic, clientID string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		BatchTimeout:           10 * time.Millisecond,
		AllowAutoTopicCreation: true,
		Transport: &kafka.Transport{
			ClientID: clientID,
		},
	}

	return &Producer{writer: w}
}

func (p *Producer) PublishCreatedOrder(ctx context.Context, event OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := strconv.FormatUint(uint64(event.UserID), 10)
	return p.Publish(ctx, key, payload)
}

func (p *Producer) PublishOrder(ctx context.Context, order database.Order) error {
	event := OrderCreatedEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Amount:    order.Amount,
		CreatedAt: order.CreatedAt,
	}

	return p.PublishCreatedOrder(ctx, event)
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
