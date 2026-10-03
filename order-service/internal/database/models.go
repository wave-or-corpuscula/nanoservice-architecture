package database

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Outbox struct {
	ID          uint      `gorm:"primaryKey"`
	EventType   string    `gorm:"not null"`
	AggregateID string    `gorm:"not null"`
	Payload     []byte    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
	PublishedAt time.Time `gorm:"type:timestamptz;default:null"`
}
