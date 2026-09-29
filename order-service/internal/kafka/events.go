package kafka

import "time"

type OrderCreatedEvent struct {
	OrderID   uint      `json:"order_id"`
	UserID    uint      `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
