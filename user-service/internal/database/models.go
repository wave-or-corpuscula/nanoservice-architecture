package database

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

type UsersPaginationResponse struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}
