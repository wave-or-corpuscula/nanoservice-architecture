package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type DLQWriter struct {
	writer *kafka.Writer
}

func NewWriter(brokers []string, topic, clientID string) *DLQWriter {
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
	return &DLQWriter{writer: w}
}

func (w *DLQWriter) WriteMessage(ctx context.Context, key []byte, value []byte) error {
	return w.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
