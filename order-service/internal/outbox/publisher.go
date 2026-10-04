package outbox

import (
	"context"
	"encoding/json"
	"log"
	"orderservice/internal/database"
	"time"
)

type Event struct {
	ID          uint      `json:"id"`
	EventType   string    `json:"event_type"`
	AggregateID string    `json:"aggregate_id"`
	Payload     []byte    `json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
}

type Database interface {
	SetEventPublished(ctx context.Context, eventID uint) error
	GetUnpublished(ctx context.Context, eventType string) ([]*database.Outbox, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, key string, value []byte) error
	PublishOrder(ctx context.Context, order database.Order) error
}

type Publisher struct {
	eventType      string
	db             Database
	eventPublisher EventPublisher
}

func New(eventType string, db Database, eventPublisher EventPublisher) *Publisher {
	return &Publisher{
		eventType:      eventType,
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (p *Publisher) Run(ctx context.Context) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := p.PublishUnpublished(ctx); err != nil {
				log.Printf("error while publishing: %v\n", err)
			}
		}
	}
}

func (p *Publisher) PublishUnpublished(ctx context.Context) error {
	events, err := p.db.GetUnpublished(ctx, p.eventType)
	if err != nil {
		return err
	}

	for _, outboxEvent := range events {
		event := Event{
			ID:          outboxEvent.ID,
			EventType:   outboxEvent.EventType,
			AggregateID: outboxEvent.AggregateID,
			Payload:     outboxEvent.Payload,
			CreatedAt:   outboxEvent.CreatedAt,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if err := p.eventPublisher.Publish(ctx, outboxEvent.AggregateID, payload); err != nil {
			return err
		} else {
			log.Printf("published: [%s] key: %s, payload: %s\n", outboxEvent.EventType, outboxEvent.AggregateID, string(outboxEvent.Payload))
		}

		if err := p.db.SetEventPublished(ctx, outboxEvent.ID); err != nil {
			return err
		}
	}
	return nil
}
