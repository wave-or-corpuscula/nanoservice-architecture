package kafka

import "time"

const EventTypeOrderCreated = "orders.created"

type Event struct {
	ID          uint      `json:"id"`
	EventType   string    `json:"event_type"`
	AggregateID string    `json:"aggregate_id"`
	Payload     []byte    `json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderCreatedEvent struct {
	OrderID   uint      `json:"order_id"`
	UserID    uint      `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
