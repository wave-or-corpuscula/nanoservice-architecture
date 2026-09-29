package handlers

import (
	"context"
	"net/http"
	"orderservice/internal/config"
	"orderservice/internal/database"
	"orderservice/internal/kafka"

	user "microservices/proto/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Database interface {
	CreateOrder(userID uint, amount float64) (*database.Order, error)
	GetUserOrders(userID uint) (*database.OrdersResponse, error)
}

type UserServiceClient interface {
	GetUser(ctx context.Context, in *user.GetUserRequest, opts ...grpc.CallOption) (*user.GetUserResponse, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, key string, value []byte) error
	PublishCreatedOrder(ctx context.Context, event kafka.OrderCreatedEvent) error
}

type Handler struct {
	db         Database
	userClient UserServiceClient
	config     config.Config
	publisher  EventPublisher
}

func NewHandler(db Database, userClient UserServiceClient, config config.Config, publisher EventPublisher) *Handler {
	return &Handler{
		db:         db,
		userClient: userClient,
		config:     config,
		publisher:  publisher,
	}
}

func (h *Handler) RegisterRouters(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/orders", h.CreateOrder)
		v1.GET("/users/:id/orders", h.GetUserOrders)
	}
}
