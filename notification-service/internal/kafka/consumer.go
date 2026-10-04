package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Database interface {
	EventExitsts(ctx context.Context, eventID uint) (bool, error)
	ProcessEvent(ctx context.Context, eventID uint, processedAt time.Time) error
}

type Consumer struct {
	db      Database
	reader  *kafka.Reader
	dlq     *DLQWriter
	retryes int
}

func NewConsumer(db Database, dlq *DLQWriter, retryes int, brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 1e6,
		},
	)

	return &Consumer{
		db:      db,
		dlq:     dlq,
		reader:  r,
		retryes: retryes,
	}
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

		var event Event
		err = json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&event)
		if err != nil {
			return err
		}

		exists, err := c.db.EventExitsts(ctx, event.ID)
		if err != nil {
			return err
		}

		if exists {
			log.Println("skipping duplicate event: ", event.ID)
			continue
		}

		var processError error
		for i := range c.retryes {
			if event.EventType == EventTypeOrderCreated {
				processError = ProcessEvent(event)
				if processError != nil {
					log.Printf("cannot process event: %v, retrying %d/%d", processError, i+1, c.retryes)
					time.Sleep(time.Second)
					continue
				} else {
					break
				}
			} else {
				log.Printf("unknown event type: %s\n", event.EventType)
			}
		}

		if processError != nil {
			if err := c.dlq.WriteMessage(ctx, msg.Key, msg.Value); err != nil {
				log.Printf("cannot publish into dlq: %v\n", err)
				return err
			} else {
				log.Println("cannot process, sent to dlq:", string(msg.Value))
			}
			continue
		}

		err = c.db.ProcessEvent(ctx, event.ID, time.Now())
		if err != nil {
			return err
		}
	}
}

func ProcessEvent(event Event) error {
	var orderEvent OrderCreatedEvent
	err := json.NewDecoder(bytes.NewReader(event.Payload)).Decode(&orderEvent)
	if err != nil {
		return err
	}

	if int(orderEvent.Amount) == 100 {
		return errors.New("invalid amount")
	}

	log.Printf("[%s] (%v), user_id: %d, amount: %f\n", event.EventType, event.CreatedAt, orderEvent.UserID, orderEvent.Amount)
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
