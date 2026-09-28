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
