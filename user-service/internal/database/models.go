package database

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	Orders    []Order   `gorm:foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

type UsersPaginationResponse struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}
