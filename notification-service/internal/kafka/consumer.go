package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 1e6,
		},
	)

	return &Consumer{reader: r}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}

		var order OrderCreatedEvent
		err = json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&order)
		if err != nil {
			return err
		}

		log.Printf("event: key=%s value=%v\n", string(msg.Key), order)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
