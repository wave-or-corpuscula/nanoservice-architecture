package database

import "time"

type ProcessedEvent struct {
	ID          uint
	ProcessedAt time.Time
}
