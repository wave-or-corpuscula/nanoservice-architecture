package service

import (
	"context"
	"encoding/json"
	"orderservice/internal/database"
	"orderservice/internal/kafka"
	"strconv"
)

type Database interface {
	WithTx(ctx context.Context, fn func(tx database.TxDatabase) error) error
	CreateOrderCtx(ctx context.Context, userID uint, amount float64) (*database.Order, error)
	CreateEvent(ctx context.Context, eventType string, key string, payload []byte) (*database.Outbox, error)
}

type OrderService struct {
	db Database
}

func NewOrderService(db Database) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, amount float64) (*database.Order, error) {
	var order *database.Order
	err := s.db.WithTx(ctx, func(tx database.TxDatabase) error {
		var err error
		order, err = tx.CreateOrderCtx(ctx, userID, amount)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(order)
		if err != nil {
			return err
		}

		if _, err := tx.CreateEvent(
			ctx,
			kafka.EventTypeOrderCreated,
			strconv.FormatUint(uint64(order.UserID), 10),
			payload,
		); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}
